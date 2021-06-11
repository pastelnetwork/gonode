package tfmodel

import (
	"context"
	"os"
	"sync"
	"testing"

	tensorflow "github.com/galeone/tensorflow/tensorflow/go"
	"github.com/pastelnetwork/gonode/probe/tfmodel/test"
	"github.com/stretchr/testify/suite"
)

type tfTestSuite struct {
	suite.Suite
	mobileNet    *TFModel
	nasNetMobile *TFModel
	modelBaseDir string
}

func TestTFModelSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(tfTestSuite))
}

func (suite *tfTestSuite) SetupSuite() {
	dir, _, err := test.UnzipModels()
	suite.modelBaseDir = dir
	suite.fatal(err)

	suite.mobileNet = &TFModel{
		Mutex: sync.Mutex{},
		Config: Config{
			path:  test.MobileNetV2,
			input: test.MobileNetV2Input,
		},
	}
	err = suite.mobileNet.Load(context.Background(), suite.modelBaseDir)
	suite.fatal(err)

	suite.nasNetMobile = &TFModel{
		Mutex: sync.Mutex{},
		Config: Config{
			path:  test.NASNetMobile,
			input: test.NASNetMobileInput,
		},
	}
	err = suite.nasNetMobile.Load(context.Background(), suite.modelBaseDir)
	suite.fatal(err)
}

func (suite *tfTestSuite) fatal(err error) {
	if err != nil {
		if suite.modelBaseDir != "" {
			os.RemoveAll(suite.modelBaseDir)
		}
		suite.T().Fatal(err)
	}
}

func (suite *tfTestSuite) TearDownSuite() {
	os.RemoveAll(suite.modelBaseDir)
}

func (suite *tfTestSuite) TestTFModel_Load() {
	type args struct {
		ctx     context.Context
		baseDir string
	}
	tests := []struct {
		name    string
		conf    Config
		args    args
		wantErr bool
	}{
		{
			name: "path is not existing",
			conf: Config{},
			args: args{
				ctx:     context.Background(),
				baseDir: "somewhere",
			},
			wantErr: true,
		},
		{
			name: "path is existing but not contains model",
			conf: Config{
				path: "somewhere",
			},
			args: args{
				ctx:     context.Background(),
				baseDir: suite.modelBaseDir,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			err := (&TFModel{
				Mutex:  sync.Mutex{},
				Config: tt.conf,
				data:   nil,
			}).Load(tt.args.ctx, tt.args.baseDir)
			if tt.wantErr {
				suite.Error(err)
			} else {
				suite.NoError(err)
			}
		})
	}
}

func (suite *tfTestSuite) TestTFModel_operation() {
	var statefulPartitionedCall = "StatefulPartitionedCall"
	type args struct {
		operation string
		index     int
	}
	type expected struct {
		dtype tensorflow.DataType
		shape []int64
	}

	tests := []struct {
		name    string
		tfModel *TFModel
		args    args
		want    expected
		wantErr bool
	}{
		{
			name:    "mobileNet operation is not existing",
			tfModel: suite.mobileNet,
			args: args{
				operation: "defcon",
				index:     0,
			},
			wantErr: true,
		},
		{
			name:    "mobileNet index is wrong",
			tfModel: suite.mobileNet,
			args: args{
				operation: statefulPartitionedCall,
				index:     1,
			},
			wantErr: true,
		},
		{
			name:    "nasNetMobile correct input operation and index",
			tfModel: suite.mobileNet,
			args: args{
				operation: test.MobileNetV2Input,
				index:     0,
			},
			want: expected{
				dtype: tensorflow.Float,
				shape: []int64{-1, -1, -1, 3},
			},
			wantErr: false,
		},
		{
			name:    "mobileNet correct operation and index",
			tfModel: suite.mobileNet,
			args: args{
				operation: statefulPartitionedCall,
				index:     0,
			},
			want: expected{
				dtype: tensorflow.Float,
				shape: []int64{-1, 1280},
			},
			wantErr: false,
		},
		{
			name:    "nasNetMobile operation is not existing",
			tfModel: suite.nasNetMobile,
			args: args{
				operation: "defcon",
				index:     0,
			},
			wantErr: true,
		},
		{
			name:    "nasNetMobile index is wrong",
			tfModel: suite.nasNetMobile,
			args: args{
				operation: statefulPartitionedCall,
				index:     1,
			},
			wantErr: true,
		},
		{
			name:    "nasNetMobile correct output operation and index",
			tfModel: suite.nasNetMobile,
			args: args{
				operation: statefulPartitionedCall,
				index:     0,
			},
			want: expected{
				dtype: tensorflow.Float,
				shape: []int64{-1, 1056},
			},
			wantErr: false,
		},
		{
			name:    "nasNetMobile correct input operation and index",
			tfModel: suite.nasNetMobile,
			args: args{
				operation: test.NASNetMobileInput,
				index:     0,
			},
			want: expected{
				dtype: tensorflow.Float,
				shape: []int64{-1, 224, 224, 3},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			output, err := tt.tfModel.operation(tt.args.operation, tt.args.index)
			if tt.wantErr {
				suite.Error(err)
				return
			}
			suite.NoError(err)
			suite.Equal(tt.args.index, output.Index)
			suite.Equal(tt.want.dtype, output.DataType())
			suite.NotNil(output.Shape())
			shape, err := output.Shape().ToSlice()
			suite.NoError(err)
			suite.Equal(tt.want.shape, shape)
		})
	}
}

func (suite *tfTestSuite) TestTFModel_Exec() {
	type args struct {
		ctx   context.Context
		value interface{}
	}

	type expected struct {
		length int
	}

	tests := []struct {
		name    string
		tfModel *TFModel
		args    args
		want    expected
		wantErr bool
	}{
		{
			name:    "mobileNet wrong input value type string",
			tfModel: suite.mobileNet,
			args: args{
				ctx:   context.Background(),
				value: "here is string",
			},
			want:    expected{},
			wantErr: true,
		},
		{
			name:    "mobileNet wrong input value type int array",
			tfModel: suite.mobileNet,
			args: args{
				ctx:   context.Background(),
				value: [1][224][224][3]int{},
			},
			want:    expected{},
			wantErr: true,
		},
		{
			name:    "mobileNet wrong input value type - incorrect array dimension",
			tfModel: suite.mobileNet,
			args: args{
				ctx:   context.Background(),
				value: [1][23][1133][33]float32{},
			},
			want:    expected{},
			wantErr: true,
		},
		{
			name:    "mobileNet correct input dimension",
			tfModel: suite.mobileNet,
			args: args{
				ctx:   context.Background(),
				value: [1][224][224][3]float32{},
			},
			want: expected{
				length: 1280,
			},
			wantErr: false,
		},
		{
			name:    "nasNetMobile wrong input value type",
			tfModel: suite.nasNetMobile,
			args: args{
				ctx:   context.Background(),
				value: "can model accept",
			},
			want:    expected{},
			wantErr: true,
		},
		{
			name:    "nasNetMobile wrong input value type - incorrect array dimension",
			tfModel: suite.nasNetMobile,
			args: args{
				ctx:   context.Background(),
				value: [1][222][224][3]float32{},
			},
			want:    expected{},
			wantErr: true,
		},
		{
			name:    "nasNetMobile correct input dimension",
			tfModel: suite.nasNetMobile,
			args: args{
				ctx:   context.Background(),
				value: [1][224][224][3]float32{},
			},
			want: expected{
				length: 1056,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			output, err := tt.tfModel.Exec(tt.args.ctx, tt.args.value)
			if tt.wantErr {
				suite.Error(err)
				return
			}
			suite.NoError(err)
			suite.Len(output, tt.want.length)
		})
	}
}
