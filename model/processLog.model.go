package model

import (
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/qq-mercantil/qq-framework-basic-golang/exceptions"
	"github.com/qq-mercantil/qq-framework-basic-golang/utils"
	"github.com/qq-mercantil/qq-framework-log-golang/logger"
	"gorm.io/gorm"
)

type ProcessLogMessageType string

const (
	INFO                ProcessLogMessageType = "INFO"
	FINISH_WITH_SUCCESS ProcessLogMessageType = "FINISH_WITH_SUCCESS"
	WARNING             ProcessLogMessageType = "WARNING"
	FINISH_WITH_ERROR   ProcessLogMessageType = "FINISH_WITH_ERROR"
	BEGINNING           ProcessLogMessageType = "BEGINNING"
	ERROR               ProcessLogMessageType = "ERROR"
)

type ProcessLogModel struct {
	ID          string                `gorm:"type:uuid;primaryKey;column:id"`
	ProcessName string                `gorm:"column:process_name;not null"`
	Props       *string               `gorm:"column:props;type:varchar"`
	Result      *string               `gorm:"column:result;type:varchar"`
	Error       *string               `gorm:"column:error;type:varchar"`
	Message     *string               `gorm:"column:message;type:varchar"`
	MessageType ProcessLogMessageType `gorm:"column:message_type;type:varchar;not null"`
	CreatedAt   time.Time             `gorm:"column:created_at;autoCreateTime"`
	StartedAt   time.Time             `gorm:"column:started_at;type:timestamp;not null"`
	FinishedAt  *time.Time            `gorm:"column:finished_at;type:timestamp"`
	ErrorClass  *string               `gorm:"column:error_class;type:varchar"`
	ErrorType   *string               `gorm:"column:error_type;type:varchar"`

	db *gorm.DB `gorm:"-"`
}

func (ProcessLogModel) TableName() string {
	return "process_log"
}

func (p *ProcessLogModel) SetDB(db *gorm.DB) {
	p.db = db
}

func (p *ProcessLogModel) Save() error {
	if p.db == nil {
		return exceptions.NewApplicationProcessException("Database connection not set")
	}
	return p.db.Save(p).Error
}

func (p *ProcessLogModel) Update(updates map[string]interface{}) error {
    if p.db == nil {
        return exceptions.NewApplicationProcessException("Database connection not set")
    }
    return p.db.Model(p).Updates(updates).Error
}

func UpdateByID(db *gorm.DB, id string, updates map[string]interface{}) error {
    return db.Model(&ProcessLogModel{}).Where("id = ?", id).Updates(updates).Error
}

type SaveErrorProps struct {
	ProcessName string
	Error       exceptions.Exception
	Message     *string
	Props       *string
}

func SaveError(db *gorm.DB, props SaveErrorProps) error {
	sLog := logger.Get()
	sLog.Error("ProcessLogModel", props.Error)

	now := time.Now()
	errorStack := props.Error.Error()
	errorClass := GetErrorClass(props.Error)
	errorType := GetErrorType(props.Error)

	instance := &ProcessLogModel{
		ID:          uuid.New().String(),
		MessageType: ERROR,
		ProcessName: props.ProcessName,
		Message:     formatPrintableMessagePtr(props.Message, 1000),
		Error:       formatPrintableMessagePtr(&errorStack, 1000),
		ErrorClass:  &errorClass,
		ErrorType:   &errorType,
		StartedAt:   now,
		FinishedAt:  &now,
		Props:       formatPrintableMessagePtr(props.Props, 1000),
		db:          db,
	}

	if err := instance.Save(); err != nil {
		sLog.Errorf("error on execute ProcessLog.SaveError(): %v", err.Error())
		return err
	}
	return nil
}

func Error(db *gorm.DB, processName string, err error, message *string) error {
	sLog := logger.Get()

	now := time.Now()
	errorStack := err.Error()
	errorClass := GetErrorClass(err)
	errorType := GetErrorType(err)

	instance := &ProcessLogModel{
		ID:          uuid.New().String(),
		MessageType: ERROR,
		ProcessName: processName,
		Message:     formatPrintableMessagePtr(message, 1000),
		Error:       formatPrintableMessagePtr(&errorStack, 1000),
		ErrorClass:  &errorClass,
		ErrorType:   &errorType,
		StartedAt:   now,
		FinishedAt:  &now,
		db:          db,
	}

	if saveErr := instance.Save(); saveErr != nil {
		sLog.Errorf("error on execute ProcessLog.Error(): %v", err.Error())
		return saveErr
	}
	return nil
}

func Info(db *gorm.DB, processName string, message string) error {
	sLog := logger.Get()

	now := time.Now()
	formattedMessage := utils.FormatPrintableMessage(message, 1000)

	instance := &ProcessLogModel{
		ID:          uuid.New().String(),
		MessageType: INFO,
		ProcessName: processName,
		Message:     &formattedMessage,
		StartedAt:   now,
		FinishedAt:  &now,
		db:          db,
	}

	if err := instance.Save(); err != nil {
		sLog.Errorf("error on execute ProcessLog.Info(): %v", err.Error())		
		return err
	}
	return nil
}

func Beginning(db *gorm.DB, processName string, props string) (*ProcessLogModel, error) {
	formattedProps := utils.FormatPrintableMessage(props, 1000)

	instance := &ProcessLogModel{
		ID:          uuid.New().String(),
		MessageType: BEGINNING,
		ProcessName: processName,
		Props:       &formattedProps,
		StartedAt:   time.Now(),
		db:          db,
	}

	if err := instance.Save(); err != nil {
		return nil, err
	}
	return instance, nil
}

func (p *ProcessLogModel) FinishWithError(err error) error {
	sLog := logger.Get()

	now := time.Now()
	errorStack := err.Error()
	errorClass := GetErrorClass(err)
	errorType := GetErrorType(err)

	p.MessageType = FINISH_WITH_ERROR
	p.ErrorClass = &errorClass
	p.ErrorType = &errorType
	p.Error = formatPrintableMessagePtr(&errorStack, 1000)
	p.FinishedAt = &now

	if saveErr := p.Save(); saveErr != nil {
		sLog.Errorf("error on execute ProcessLog.FinishWithError(): %v", err.Error())
		return saveErr
	}
	return nil
}

func (p *ProcessLogModel) FinishWithSuccess(result string) error {
	sLog := logger.Get()

	now := time.Now()
	formattedResult := utils.FormatPrintableMessage(result, 1000)

	p.MessageType = FINISH_WITH_SUCCESS
	p.Result = &formattedResult
	p.FinishedAt = &now

	if err := p.Save(); err != nil {
		sLog.Errorf("error on execute ProcessLog.FinishWithSuccess(): %v", err.Error())
		return err
	}
	return nil
}

func formatPrintableMessagePtr(message *string, maxLength int) *string {
	if message == nil {
		return nil
	}
	formatted := utils.FormatPrintableMessage(*message, maxLength)
	return &formatted
}

func GetErrorClass(err error) string {
	typeName := reflect.TypeOf(err).String()
	if typeName == "" {
		return "UnknownError"
	}

	return typeName
}

func GetErrorType(err error) string {
	if _, ok := err.(*exceptions.Warning); ok {
		return "Warning"
	}
	return "Error"
}

