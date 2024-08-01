package telegram

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/alexKudryavtsev-web/grace_links_tg_bot/lib/e"
)

const (
	getUpdatesMethod  = "getUpdates"
	sendMessageMethod = "sendMessage"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

func New(host string, token string) Client {
	return Client{
		host:     host,
		basePath: newBaseToken(token),
		client:   http.Client{},
	}
}

func newBaseToken(token string) string {
	return "bot" + token
}

func (c *Client) Updates(offset int, limit int) (data []Update, err error) {
	defer func() { err = e.WrapIfErr("telegram.Updates", err) }()

	q := url.Values{}

	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	req, err := c.doRequest(getUpdatesMethod, q)

	if err != nil {
		return nil, err
	}

	var res UpdatesResponse

	if err := json.Unmarshal(req, &res); err != nil {
		return nil, err
	}

	return res.Result, nil
}

func (c *Client) SendMessage(chatID int, text string) error {
	q := url.Values{}

	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)

	_, err := c.doRequest(sendMessageMethod, q)
	if err != nil {
		return e.Wrap("telegram.SendMessage: can't send messsage", err)
	}

	return nil
}

func (c *Client) doRequest(method string, query url.Values) (data []byte, err error) {
	defer func() { err = e.WrapIfErr("telegram.doRequest", err) }()

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)

	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)

	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return body, err
	}

	return body, nil
}
