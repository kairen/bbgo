package slacknotifier

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

type SlackAttachmentCreator interface {
	SlackAttachment() slack.Attachment
}

type Notifier struct {
	client  *slack.Client
	channel string
}

type NotifyOption func(notifier *Notifier)

func New(token, channel string, options ...NotifyOption) *Notifier {
	// var client = slack.New(token, slack.OptionDebug(true))
	var client = slack.New(token)

	notifier := &Notifier{
		channel: channel,
		client:  client,
	}

	for _, o := range options {
		o(notifier)
	}

	return notifier
}

func (n *Notifier) Notify(obj interface{}, args ...interface{}) {
	n.NotifyTo(n.channel, obj, args...)
}

func filterSlackAttachments(args []interface{}) (slackAttachments []slack.Attachment, pureArgs []interface{}) {
	var firstAttachmentOffset = -1
	for idx, arg := range args {
		switch a := arg.(type) {

		// concrete type assert first
		case slack.Attachment:
			if firstAttachmentOffset == -1 {
				firstAttachmentOffset = idx
			}

			slackAttachments = append(slackAttachments, a)

		case SlackAttachmentCreator:
			if firstAttachmentOffset == -1 {
				firstAttachmentOffset = idx
			}

			slackAttachments = append(slackAttachments, a.SlackAttachment())

		}
	}

	pureArgs = args
	if firstAttachmentOffset > -1 {
		pureArgs = args[:firstAttachmentOffset]
	}
	return
}

func (n *Notifier) NotifyTo(channel string, obj interface{}, args ...interface{}) {
	if len(channel) == 0 {
		channel = n.channel
	}

	slackAttachments, pureArgs := filterSlackAttachments(args)

	var opts []slack.MsgOption

	switch a := obj.(type) {
	case string:
		opts = append(opts, slack.MsgOptionText(fmt.Sprintf(a, pureArgs...), true),
			slack.MsgOptionAttachments(slackAttachments...))

	case slack.Attachment:
		opts = append(opts, slack.MsgOptionAttachments(append([]slack.Attachment{a}, slackAttachments...)...))

	case SlackAttachmentCreator:
		// convert object to slack attachment (if supported)
		opts = append(opts, slack.MsgOptionAttachments(append([]slack.Attachment{a.SlackAttachment()}, slackAttachments...)...))

	default:
		log.Errorf("slack notifier error, unsupported object: %T %+v", a, a)

	}

	go func() {
		_, _, err := n.client.PostMessageContext(context.Background(), channel, opts...)
		if err != nil {
			log.WithError(err).
				WithField("channel", channel).
				Errorf("slack api error: %s", err.Error())
		}
	}()

	return
}

/*
func (n *Notifier) NotifyTrade(trade *types.Trade) {
	_, _, err := n.client.PostMessageContext(context.Background(), n.TradeChannel,
		slack.MsgOptionText(util.Render(`:handshake: {{ .Symbol }} {{ .Side }} Trade Execution @ {{ .Price  }}`, trade), true),
		slack.MsgOptionAttachments(trade.SlackAttachment()))

	if err != nil {
		logrus.WithError(err).Error("slack send error")
	}
}
*/

/*
func (n *Notifier) NotifyPnL(report *pnl.AverageCostPnlReport) {
	attachment := report.SlackAttachment()

	_, _, err := n.client.PostMessageContext(context.Background(), n.PnlChannel,
		slack.MsgOptionText(util.Render(
			`:heavy_dollar_sign: Here is your *{{ .symbol }}* PnL report collected since *{{ .startTime }}*`,
			map[string]interface{}{
				"symbol":    report.Symbol,
				"startTime": report.StartTime.Format(time.RFC822),
			}), true),
		slack.MsgOptionAttachments(attachment))

	if err != nil {
		logrus.WithError(err).Errorf("slack send error")
	}
}
*/
