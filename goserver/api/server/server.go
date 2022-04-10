package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"goserver/api/types/files"
	pb "goserver/api/types/protos"
	"goserver/configs/exploredirs"
	"goserver/pkg/dbclient"
	"goserver/pkg/filehandler"

	"google.golang.org/grpc"
)

type raxFileServer struct {
	pb.UnimplementedRaxFileServer
}

func localFileInfo2RPC(localFileInfo *files.FileInfo) *pb.FileInfo {
	if localFileInfo == nil {
		return nil
	}

	return &pb.FileInfo{
		ID:           *localFileInfo.ID,
		FilePath:     *localFileInfo.FilePath,
		FileName:     *localFileInfo.FileName,
		ExtName:      *localFileInfo.ExtName,
		CreationTime: *localFileInfo.CreationTime,
		Size:         *localFileInfo.CreationTime,
		Extra:        *localFileInfo.Extra,
	}
}
func (r *raxFileServer) SearchVideos(ctx context.Context, in *pb.SearchVideosRequest) (*pb.SearchVideosResponse, error) {
	fileInfos, total := dbclient.SearchVideos(in.Patterns, int(in.PageNo), int(in.PageSize))

	pbFileInfos := []*pb.FileInfo{}
	for _, finfo := range fileInfos {
		pbFileInfos = append(pbFileInfos, localFileInfo2RPC(finfo))
	}
	res := pb.SearchVideosResponse{
		FileInfos: pbFileInfos,
		Total:     int32(total),
	}

	return &res, nil
}

func (r *raxFileServer) SearchVideoByID(ctx context.Context, in *pb.SearchVideoByIDRequest) (*pb.SearchVideoByIDResponse, error) {
	fInfo := dbclient.SearchVideoByID(in.ID)

	res := pb.SearchVideoByIDResponse{
		FileInfo: localFileInfo2RPC(fInfo),
	}

	return &res, nil
}
func (r *raxFileServer) RescanFilesAndUpdateDB(ctx context.Context, in *pb.RescanFilesAndUpdateDBRequest) (*pb.RescanFilesAndUpdateDBResponse, error) {
	// search all files
	scanT := time.Now()
	fileInfos := filehandler.ExploreFilesMulti(exploredirs.Dirs)
	fmt.Println("scan", time.Since(scanT))

	{
		insertT := time.Now()
		defer func() {
			fmt.Println("insert", time.Since(insertT))
		}()
	}

	if err := dbclient.UpdateVideosSoft(fileInfos); err == nil {
		return &pb.RescanFilesAndUpdateDBResponse{
			Success: true,
		}, nil
	} else {
		return &pb.RescanFilesAndUpdateDBResponse{
			Success: false,
		}, err
	}
}

func main() {
	lis, err := net.Listen("tcp", "127.0.0.1:6969")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterRaxFileServer(s, &raxFileServer{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
