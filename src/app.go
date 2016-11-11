package kmstool

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"gopkg.in/urfave/cli.v1"
	"strings"
)

var parserRegex string
var masterKey string
var awsRegion string
var isTest bool

type App struct {
	*cli.App
	Client *kms.KMS
	Item   KMSItem
}

func (s *App) CommandDecrypt(c *cli.Context) error {
	err := s.Item.RunOnParsed(func(parsedHit string, parsedContent string) (content string, err error) {
		return Decrypt(s.Client, parsedContent)
	})
	if err != nil {
		return err
	}

	s.Item.Replace(func(content string, parsedHit string, parsedContent string) (newContent string, err error) {
		return strings.Replace(content, parsedHit, parsedContent, -1), nil
	})

	_, err = s.Writer.Write([]byte(s.Item.Content()))

	return err
}

func (s *App) CommandEncrypt(c *cli.Context) error {
	if len(masterKey) == 0 {
		return errors.New("You must specify master key to encrypt")
	}

	err := s.Item.RunOnParsed(func(parsedHit string, parsedContent string) (content string, err error) {
		return Encrypt(s.Client, masterKey, parsedContent)
	})
	if err != nil {
		return err
	}

	s.Item.Replace(func(content string, parsedHit string, parsedContent string) (newContent string, err error) {
		return strings.Replace(content, parsedHit, parsedContent, -1), nil
	})

	s.Writer.Write([]byte(s.Item.Content()))
	return nil
}

func (s *App) Before(c *cli.Context) error {
	var err error
	var awsSession *session.Session

	item := NewKMSItem(c.Args().First())

	s.Item = *item
	err = s.Item.Parse(NewRegexpStringParser(parserRegex))

	if err != nil {
		return err
	}

	if s.Client == nil {
		awsSession, err = session.NewSession()
		if err != nil {
			return err
		}
		s.Client = kms.New(awsSession, &aws.Config{Region: &awsRegion})
	}

	return nil
}

func NewApp() *App {
	app := &App{
		cli.NewApp(),
		nil,
		KMSItem{},
	}
	app.Name = "kmstool"
	app.Usage = "Encrypt/decrypt using AWS KMS. " +
		"Login credentials are read from env {AWS_ACCESS_KEY_ID, AWS_SECRET_KEY, AWS_SESSION_TOKEN} or by specifying " +
		"AWS_SHARED_CREDENTIALS_FILE=~/.aws/config"

	app.Version = "0.1.0"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "regexp, r",
			Value:       `(?P<hit>kms:(?P<key>.+))(?:\n|"|'|$)`,
			Usage:       "Regexp with 'hit' and 'key' groups. 'hit' is used for replacing, 'key' used for commands",
			Destination: &parserRegex,
		},
		cli.StringFlag{
			Name:        "cmk, k",
			Usage:       "KMS master key",
			Destination: &masterKey,
		},
		cli.StringFlag{
			Name:        "aws-region, region",
			Usage:       "AWS Region",
			Destination: &awsRegion,
			EnvVar:      "AWS_REGION",
		},
		cli.BoolFlag{
			Name:        "test",
			Usage:       "",
			Destination: &isTest,
			EnvVar:      "KMSTOOL_TEST",
			Hidden:      true,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "encrypt",
			Aliases: []string{"e"},
			Usage:   "Encrypt and replace multiple 'key' in 'content'. You need to spefiy aws master key with '--cmk' option",
			Action:  app.CommandEncrypt,
			Before:  app.Before,
		},
		{
			Name:    "decrypt",
			Aliases: []string{"d"},
			Usage:   "Decrypt and replace multiple 'key' in 'content'",
			Action:  app.CommandDecrypt,
			Before:  app.Before,
		},
	}

	return app
}
