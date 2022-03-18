package dc_test

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/cyrusaf/dc"
	"github.com/cyrusaf/dc/parser"
	"github.com/cyrusaf/dc/source"
)

func ExampleConfigLocal() {
	type FeatureConfig struct {
		EnableMyFeature bool
	}

	ctx := context.Background()

	localSource := source.Local[FeatureConfig]{
		Data: FeatureConfig{
			EnableMyFeature: true,
		},
	}

	featureConfig, err := dc.NewConfig[FeatureConfig](ctx, localSource, FeatureConfig{})
	if err != nil {
		// It is good practice to log on initialize failure, but continue. We
		// don't want a source outage to cause our app to fail to start if we
		// can rely on sane defaults for the config.
		fmt.Printf("error initializing feature config: %v\n", err)
	}

	if featureConfig.Get().EnableMyFeature {
		fmt.Println("MyFeature enabled!")
		return
	}
	fmt.Println("MyFeature disabled!")

	// Output: MyFeature enabled!
}

func ExampleConfigS3() {
	type FeatureConfig struct {
		EnableMyFeature bool `json:"enable_my_feature"`
	}

	ctx := context.Background()

	sess := session.Must(session.NewSession())
	s3Downloader := s3manager.NewDownloader(sess)
	s3Source := source.S3[FeatureConfig]{
		BucketName:   "my_bucket",
		Key:          "configs/my_config.json",
		S3Downloader: s3Downloader,
		Parser:       parser.JSON[FeatureConfig]{},
	}

	featureConfig, err := dc.NewConfig[FeatureConfig](ctx, s3Source, FeatureConfig{})
	if err != nil {
		// It is good practice to log on initialize failure, but continue. We
		// don't want a source outage to cause our app to fail to start if we
		// can rely on sane defaults for the config.
		fmt.Printf("error initializing feature config: %v\n", err)
	}

	poller := dc.NewPoller(featureConfig, time.Second*5, func(err error) {
		fmt.Printf("error polling feature config: %v\n", err)
	})
	go poller.Poll(ctx)
	defer poller.Shutdown()

	if featureConfig.Get().EnableMyFeature {
		fmt.Println("MyFeature enabled!")
		return
	}
	fmt.Println("MyFeature disabled!")

	// Output: MyFeature disabled!
}
