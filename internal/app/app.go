package app

import (
	"context"
	"path/filepath"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/core-go/io/importer"
	"github.com/core-go/io/reader"
	"github.com/core-go/io/transform"
	v "github.com/core-go/io/validator"
	"github.com/core-go/log"
	w "github.com/core-go/mongo/writer"
)

type ApplicationContext struct {
	Import func(ctx context.Context) (int, int, error)
}

func NewApp(ctx context.Context, cfg Config) (*ApplicationContext, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Mongo.Uri))
	db := client.Database(cfg.Mongo.Database)
	if err != nil {
		return nil, err
	}
	fileType := reader.DelimiterType
	filename := ""
	if fileType == reader.DelimiterType {
		filename = "delimiter.csv"
	} else {
		filename = "fixedlength.csv"
	}
	generateFileName := func() string {
		fullPath := filepath.Join("export", filename)
		return fullPath
	}
	reader, err := reader.NewDelimiterFileReader(generateFileName)
	if err != nil {
		return nil, err
	}
	mp := map[string]interface{}{
		"app": "import users",
		"env": "dev",
	}
	transformer, err := transform.NewDelimiterTransformer[User](",")
	if err != nil {
		return nil, err
	}
	validator, err := v.NewValidator[*User]()
	if err != nil {
		return nil, err
	}
	errorHandler := importer.NewErrorHandler[*User](log.ErrorFields, "fileName", "lineNo", mp)
	writer := w.NewMongoWriter[*User](db, "userimport")
	importer := importer.NewImporter[User](reader.Read, transformer.Transform, validator.Validate, errorHandler.HandleError, errorHandler.HandleException, filename, writer.Write)
	return &ApplicationContext{Import: importer.Import}, nil
}

type User struct {
	Id          string     `json:"id" gorm:"column:id;primary_key" bson:"_id" format:"%011s" length:"11" dynamodbav:"id" firestore:"id" validate:"required,max=40"`
	Username    string     `json:"username" gorm:"column:username" bson:"username" length:"10" dynamodbav:"username" firestore:"username" validate:"required,username,max=100"`
	Email       string     `json:"email" gorm:"column:email" bson:"email" dynamodbav:"email" firestore:"email" length:"31" validate:"email,max=100"`
	Phone       string     `json:"phone" gorm:"column:phone" bson:"phone" dynamodbav:"phone" firestore:"phone" length:"20" validate:"required,max=18"`
	Status      bool       `json:"status" gorm:"column:status" true:"1" false:"0" bson:"status" dynamodbav:"status" format:"%5s" length:"5" firestore:"status"`
	CreatedDate *time.Time `json:"createdDate" gorm:"column:createdDate" bson:"createdDate" length:"10" format:"dateFormat:2006-01-02" dynamodbav:"createdDate" firestore:"createdDate" validate:"required"`
	Test        string     `json:"test" gorm:"-" bson:"-" dynamodbav:"-" firestore:"-" length:"0" format:"-"`
}
