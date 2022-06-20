package resizer

import (
	"context"
	"image"
	"image/jpeg"
	"image/png"
	"os"

	"github.com/palantir/stacktrace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"

	"github.com/Raphy42/weekend/core"
	"github.com/Raphy42/weekend/core/errors"
	"github.com/Raphy42/weekend/core/logger"
	"github.com/Raphy42/weekend/pkg/slice"

	"golang.org/x/image/draw"
)

type Input struct {
	InputPath  string `json:"inputPath" validate:"required,file"`
	OutputPath string `json:"outputPath" validate:"required,file"`
	OutMime    string `json:"outMime" validate:"required,oneof=jpeg jpg png"`
	NewWidth   int    `json:"newWidth" validate:"gt=0"`
	NewHeight  int    `json:"newHeight" validate:"gt=0"`
	// unused because catmull-rom is selected by default in order to have actual load
	//Efficiency float32 `json:"efficiency"`
}

type Output struct {
	Size int `json:"size"`
}

func ResizeImage(ctx context.Context, args Input) (*Output, error) {
	ctx, span := otel.Tracer(core.Name()).Start(ctx, "ResizeImage")
	dims := slice.New(args.NewWidth, args.NewHeight)
	span.SetAttributes(
		attribute.String("wk.task.args.input", args.InputPath),
		attribute.String("wk.task.args.output", args.InputPath),
		attribute.String("wk.task.args.outMime", args.OutMime),
		attribute.IntSlice("wk.task.args.newDims", dims),
	)
	defer span.End()

	log := logger.FromContext(ctx)
	log.Info("resizing image",
		zap.String("filename", args.InputPath),
		zap.Ints("dimensions", dims),
	)

	file, err := os.Open(args.InputPath)
	defer func() {
		if err = file.Close(); err != nil {
			log.Error("io error", zap.Error(err))
		}
	}()

	if err != nil {
		return nil, stacktrace.Propagate(err, "could not open input image: '%s'", args.InputPath)
	}

	output, err := os.Create(args.OutputPath)
	defer func() {
		if err = output.Close(); err != nil {
			log.Error("io error", zap.Error(err))
		}
	}()

	src, _, err := image.Decode(file)
	if err != nil {
		return nil, stacktrace.Propagate(err, "image decoding failed")
	}

	newBounds := image.Rect(0, 0, args.NewWidth, args.NewHeight)
	dst := image.NewRGBA(newBounds)
	draw.CatmullRom.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)

	switch args.OutMime {
	case "jpeg", "jpg":
		if err := jpeg.Encode(output, dst, &jpeg.Options{Quality: 75}); err != nil {
			return nil, err
		}
	case "png":
		if err = png.Encode(output, dst); err != nil {
			return nil, err
		}
	default:
		return nil, stacktrace.NewErrorWithCode(
			errors.EUnreachable,
			"unhandled output mime type: '%s'", args.OutMime,
		)
	}

	return &Output{Size: len(dst.Pix)}, nil
}
