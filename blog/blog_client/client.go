package main

import (
	"context"
	"fmt"
	"github.com/jirawan-chuapradit/grpc-go-course/blog/blogpb"
	"google.golang.org/grpc"
	"log"
)

func main() {

	fmt.Println("Blog Client")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)

	// create Blog
	fmt.Println("creating the blog")
	blog := &blogpb.Blog{
		AuthorId: "Stephane",
		Title: "My First Blog",
		Content: "Content of the first blog",
	}
	createBlogRes, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err!=nil{
		log.Fatalf("unexpected error: %v\n", err)
	}
	fmt.Printf("Blog has been created: %v\n", createBlogRes)


	//read Blog
	fmt.Println("reading the blog")
	_, err = c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{
		BlogId: "23dsfse",
	})
	if err != nil {
		fmt.Printf("error happened while reading: %v\n", err)
	}

	fmt.Println()
	fmt.Println()
	res2, err := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{
		BlogId: createBlogRes.GetBlog().GetId(),
	})
	if err != nil {
		fmt.Sprintf("error happened while reading: %v\n", err)
	}
	fmt.Printf("blog was read: %v", res2)
}
