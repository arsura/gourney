package repository_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/arsura/gourney/config"
	adapter "github.com/arsura/gourney/pkg/adapters"
	"github.com/arsura/gourney/pkg/adapters/mocks"
	model "github.com/arsura/gourney/pkg/models/mongodb"
	repository "github.com/arsura/gourney/pkg/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type PostRepositoryTestSuite struct {
	suite.Suite
	mockPostCollection *mocks.MongoCollectionProvider
	postRepository     repository.PostRepositoryProvider
}

func (suite *PostRepositoryTestSuite) SetupTest() {
	suite.mockPostCollection = new(mocks.MongoCollectionProvider)
	mockMongoCollections := &adapter.MongoCollections{
		BlogDatabase: &adapter.BlogDatabase{
			PostCollection: suite.mockPostCollection,
		},
	}
	suite.postRepository = repository.NewPostRepository(mockMongoCollections, &zap.SugaredLogger{}, &config.Config{})
}

func (suite *PostRepositoryTestSuite) TestCreatePostSuccess() {
	now := time.Now()
	mockPost := &model.Post{
		Title:     "TestCreate",
		Content:   "TestCreatePostSuccess",
		CreatedAt: now,
		UpdatedAt: now,
	}
	mockId := primitive.NewObjectID()
	suite.mockPostCollection.On("InsertOne", mock.Anything, mockPost).Return(&mongo.InsertOneResult{InsertedID: mockId}, nil)
	result, err := suite.postRepository.CreatePost(context.Background(), mockPost)
	suite.mockPostCollection.AssertNumberOfCalls(suite.T(), "InsertOne", 1)
	assert.Equal(suite.T(), result, &mockId)
	assert.Nil(suite.T(), err)
}

func (suite *PostRepositoryTestSuite) TestCreatePostFailed() {
	now := time.Now()
	mockPost := &model.Post{
		Title:     "TestCreate",
		Content:   "TestCreatePostSuccess",
		CreatedAt: now,
		UpdatedAt: now,
	}
	suite.mockPostCollection.On("InsertOne", mock.Anything, mockPost).Return(nil, errors.New("failed to insert one"))
	result, err := suite.postRepository.CreatePost(context.Background(), mockPost)
	suite.mockPostCollection.AssertNumberOfCalls(suite.T(), "InsertOne", 1)
	assert.Nil(suite.T(), result)
	assert.NotNil(suite.T(), err)
}

func TestPostRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(PostRepositoryTestSuite))
}
