package repositories

import (
	"database/sql"
	"sample-go-rest/entities"
	"github.com/Masterminds/squirrel"
	"errors"
	"time"
)

type ITopicRepository interface {
	GetList(params map[string][]string) ([]entities.Topic, error)
	GetItemById(id int) (*entities.Topic, error)
	InsertItem(Topic entities.Topic) error
	UpdateItem(id int, Topic entities.Topic) error
	DeleteItem(id int) error

	GetListByNewsId(newsId int, params map[string][]string) ([]entities.Topic, error)
	GetItemByNewsId(newsId int, topicId int) (*entities.Topic, error)
}

type TopicRepostory struct {
	db *sql.DB
}

func NewTopicRepository(db *sql.DB) ITopicRepository {
	return &TopicRepostory{db: db}
}

func (repository *TopicRepostory) GetList(params map[string][]string) ([]entities.Topic, error) {
	builder := squirrel.Select("*").From("topic")
	builder, err := applyFilterAndPageSize(builder, params)
	if err != nil {
		return nil, err
	}
	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}
	result, err := repository.db.Query(sql, args...)
	if err != nil {
		return nil, err
	}

	list := []entities.Topic{}

	for result.Next() {
		topic := scanTopic(result)

		list = append(list, topic)
	}

	return list, nil
}

func (repository *TopicRepostory) GetItemById(id int) (*entities.Topic, error) {
	sql, args, err := squirrel.Select("*").From("topic").Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, err
	}
	result, err := repository.db.Query(sql, args...)
	if err != nil {
		return nil, err
	}

	for result.Next() {
		topic := scanTopic(result)

		return &topic, nil
	}

	return nil, errors.New("item not found")
}

func (repository *TopicRepostory) InsertItem(topic entities.Topic) error {
	sql, args, err := squirrel.Insert("topic").Columns("title").Values(topic.Title).ToSql()
	if err != nil {
		return err
	}

	_, err = repository.db.Exec(sql, args...)
	return err
}

func (repository *TopicRepostory) UpdateItem(id int, topic entities.Topic) error {
	sql, args, err := squirrel.Update("topic").
		SetMap(map[string]interface{}{"title": topic.Title, "updated_at": time.Now()}).
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = repository.db.Exec(sql, args...)
	return err
}

func (repository *TopicRepostory) DeleteItem(id int) error {
	sql, args, err := squirrel.Delete("topic").Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		return err
	}

	_, err = repository.db.Exec(sql, args...)
	return err
}

func (repository *TopicRepostory) GetListByNewsId(newsId int, params map[string][]string) ([]entities.Topic, error) {
	builder := squirrel.Select("topic.*").From("topic").
		Join("news_topic ON news_topic.topic_id = topic.id").
		Where(squirrel.Eq{"news_topic.news_id": newsId})

	builder, err := applyFilterAndPageSize(builder, params)

	if err != nil {
		return nil, err
	}

	sql, args, err := builder.ToSql()

	if err != nil {
		return nil, err
	}
	result, err := repository.db.Query(sql, args...)
	if err != nil {
		return nil, err
	}

	list := []entities.Topic{}

	for result.Next() {
		topic := scanTopic(result)

		list = append(list, topic)
	}

	return list, nil
}

func (repository *TopicRepostory) GetItemByNewsId(newsId int, topicId int) (*entities.Topic, error) {
	sql, args, err := squirrel.Select("topic.*").From("topic").
		Join("news_topic ON news_topic.topic_id = topic.id").
		Where(squirrel.Eq{"news_topic.news_id": newsId, "topic.id": topicId}).
		ToSql()
	if err != nil {
		return nil, err
	}
	result, err := repository.db.Query(sql, args...)
	if err != nil {
		return nil, err
	}

	for result.Next() {
		topic := scanTopic(result)

		return &topic, nil
	}

	return nil, errors.New("item not found")
}

func scanTopic(result *sql.Rows) entities.Topic {
	var topic entities.Topic
	result.Scan(&topic.Id, &topic.Title, &topic.CreatedAt, &topic.UpdatedAt)
	return topic
}
