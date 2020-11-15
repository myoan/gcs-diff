package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"sync"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func compare(ctx context.Context, bname string, bh *storage.BucketHandle, src, dst, path string) {
	srcPath := filepath.Join(src, path)
	dstPath := filepath.Join(dst, path)

	sAttr, err := bh.Object(srcPath).Attrs(ctx)
	if err != nil {
		fmt.Printf("%v: %s\n", err, filepath.Join(fmt.Sprintf("gs://%s", bname), srcPath))
		return
	}
	dAttr, err := bh.Object(dstPath).Attrs(ctx)
	if err != nil {
		fmt.Printf("%v: %s\n", err, filepath.Join(fmt.Sprintf("gs://%s", bname), dstPath))
		return
	}
	if sAttr.CRC32C != dAttr.CRC32C {
		fmt.Printf("storage: Object doesn't match: %s\n", filepath.Join(fmt.Sprintf("gs://%s", bname), srcPath))
		return
	}
	return
}

func sliceUnique(target []string) (unique []string) {
	m := map[string]bool{}

	for _, v := range target {
		if !m[v] {
			m[v] = true
			unique = append(unique, v)
		}
	}

	return unique
}

func walk(ctx context.Context, bh *storage.BucketHandle, src, dst string) []string {
	var names []string
	for _, path := range []string{src, dst} {
		query := &storage.Query{Prefix: path}

		it := bh.Objects(ctx, query)
		for {
			attrs, err := it.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			names = append(names, strings.Replace(attrs.Name, path, "", 1))
		}
	}

	return sliceUnique(names)
}

func main() {
	var (
		bn   = flag.String("b", "", "backetname")
		cr   = flag.String("cred", "", "credential path")
		src  = flag.String("src", "", "src path to compare")
		dst  = flag.String("dst", "", "dst path to compare")
		conc = flag.Int("conc", 4, "upload cuncurrency")
	)
	flag.Parse()

	if len(*bn) == 0 {
		panic("err: undefined bucket name")
	}

	if len(*cr) == 0 {
		panic("err: undefined credential path")
	}

	if len(*src) == 0 {
		panic("err: undefined source hash")
	}

	if len(*dst) == 0 {
		panic("err: undefined destination hash")
	}

	fmt.Printf("---------------------------------------------\n")
	fmt.Printf("bucket:      %s\n", *bn)
	fmt.Printf("credential:  %s\n", *cr)
	fmt.Printf("source:      %s\n", *src)
	fmt.Printf("destination: %s\n", *dst)
	fmt.Printf("---------------------------------------------\n\n")

	ctx := context.Background()

	client, err := storage.NewClient(ctx, option.WithCredentialsFile(*cr))
	if err != nil {
		panic("err: failed to create gcs client")
	}

	b := client.Bucket(*bn)
	if _, err = b.Attrs(ctx); err != nil {
		fmt.Println(err)
		panic("bucket not found")
	}

	limit := make(chan struct{}, *conc)
	var wg sync.WaitGroup
	for _, f := range walk(ctx, b, *src, *dst) {
		wg.Add(1)
		go func(ctx context.Context, bname string, bh *storage.BucketHandle, src, dst, path string) {
			limit <- struct{}{}
			defer wg.Done()
			compare(ctx, bname, bh, src, dst, path)
			<-limit
		}(ctx, *bn, b, *src, *dst, f)
	}
	wg.Wait()
}
