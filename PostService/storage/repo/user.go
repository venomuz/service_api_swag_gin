package repo

import (
	pb "github.com/venomuz/service_api_swag_gin/PostService/genproto"
)

//PostStorageI ...
type PostStorageI interface {
	PostCreate(*pb.Post) (*pb.Post, error)
	PostGetByID(ID string) (*pb.Post, error)
	PostDeleteByID(ID string) (*pb.OkBOOL, error)
	PostGetAllPosts(ID string) (*pb.AllPost, error)
}
