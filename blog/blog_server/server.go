package main

import (
	"context"
	"fmt"
	"github.com/jirawan-chuapradit/grpc-go-course/blog/blogpb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

var collection *mongo.Collection

type server struct{}

func (s server) DeleteBlog(ctx context.Context, request *blogpb.DeleteBlogRequest) (*blogpb.DeleteBlogResponse, error) {
	fmt.Println("delete blog request ...")
	oid, err := primitive.ObjectIDFromHex(request.GetBlog().GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "cannot parse ID")
	}

	filter := bson.D{{"_id", oid}}

	res, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("cannot delete object in MongoDB: %v", err))
	}

	if res.DeletedCount == 0 {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("cannot delete object in MongoDB: %v", err))
	}

	return &blogpb.DeleteBlogResponse{
		BlogId: request.GetBlog().GetId(),
	}, nil
}

func (s server) UpdateBlog(ctx context.Context, request *blogpb.UpdateBlogRequest) (*blogpb.UpdateBlogResponse, error) {
	fmt.Println("update blog request ...")

	blog := request.GetBlog()
	oid, err := primitive.ObjectIDFromHex(blog.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "cannot parse ID")
	}

	data := &blogItem{}
	filter := bson.D{{"_id", oid}}

	err = collection.FindOne(context.Background(),filter).Decode(&data)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("cannot find blog with specified ID: %v", err))
	}

	// we update our internal struct
	data.AuthorID = blog.GetAuthorId()
	data.Content = blog.GetContent()
	data.Title = blog.GetTitle()

	_, err = collection.ReplaceOne(context.Background(), filter,data)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("cannot update object in mongoDB: %V", err))
	}
	return &blogpb.UpdateBlogResponse{
		Blog: dataToBlogPb(*data),
	},nil
}

func dataToBlogPb (data blogItem) *blogpb.Blog{
	return &blogpb.Blog{
		Id: data.ID.Hex(),
		AuthorId: data.AuthorID,
		Content: data.Content,
		Title: data.Title,
	}
}

func (s server) ReadBlog(ctx context.Context, request *blogpb.ReadBlogRequest) (*blogpb.ReadBlogResponse, error) {
	fmt.Println("read blog request ...")

	blogID := request.GetBlogId()

	oid , err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "cannot parse ID")
	}

	// create an empty struct
	data := &blogItem{}

	filter := bson.D{{"_id", oid}}
	err = collection.FindOne(context.Background(), filter).Decode(&data)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("cannot find blog with specified ID: %v", err))
	}

	return &blogpb.ReadBlogResponse{
		Blog: dataToBlogPb(*data),
	},nil
}

func (s server) CreateBlog(ctx context.Context, request *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {
	fmt.Println("create blog request ...")
	blog := request.GetBlog()

	data := blogItem{
		AuthorID: blog.GetAuthorId(),
		Title: blog.GetTitle(),
		Content: blog.GetContent(),
	}

	res, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err),
			)
	}
	old, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("cannot convert to OID"))
	}
	return &blogpb.CreateBlogResponse{
		Blog: &blogpb.Blog{
			Id: old.Hex(),
			AuthorId: blog.GetAuthorId(),
			Title: blog.GetTitle(),
			Content: blog.GetContent(),
		},
	}, nil
}

type blogItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string             `bson:"author_id"`
	Content  string             `bson:"content"`
	Title    string             `bson:"title"`
}

func main() {
	// fi we crash the go code, we get the file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("connecting to MongoDB")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

	collection = client.Database("mydb").Collection("blog")

	fmt.Println("Blog Service Started")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	blogpb.RegisterBlogServiceServer(s, &server{})

	go func() {
		fmt.Println("starting server ...")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// wait for control c to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Block until a signal is received
	<-ch
	fmt.Println("Stopping the server")
	s.Stop()
	fmt.Println("Closing the listener")
	lis.Close()
	fmt.Println("closing mongo db connection")
	client.Disconnect(context.TODO())
	fmt.Println("End of Program")
}
