package sendowl

import (
	"net/http"
	"net/url"
	"sync"

	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/urlfetch"
)

var theClient struct {
	sync.RWMutex
	datastoreClient
}

type datastoreClient struct {
	fetched  bool
	Key      string
	Secret   string
	Endpoint string
	Client   `datastore:"-"`
}

const (
	clientKind       = "SendowlClient"
	defaultClientKey = "default"
)

func (c datastoreClient) key(ctx context.Context) *datastore.Key {
	return datastore.NewKey(ctx, clientKind, defaultClientKey, 0, nil)
}

func newURLFetchTransportFunc(ctx context.Context) TransportFunc {
	return func(_ context.Context) http.RoundTripper {
		return &urlfetch.Transport{Context: ctx}
	}
}

func DefaultClient(ctx context.Context) (Client, error) {
	theClient.RLock()
	c := theClient.datastoreClient
	theClient.RUnlock()
	if c.fetched {
		return *c.Client.WithTransportFunc(newURLFetchTransportFunc(ctx)), nil
	}
	theClient.Lock()
	defer theClient.Unlock()
	if theClient.datastoreClient.fetched {
		return theClient.datastoreClient.Client, nil
	}
	key := theClient.datastoreClient.key(ctx)
	if err := datastore.Get(ctx, key, &theClient.datastoreClient); err != nil {
		if err != datastore.ErrNoSuchEntity {
			return Client{}, err
		}
		// Fill with a placeholder making it easy to change manually.
		theClient.datastoreClient = placeholderDatastoreClient(ctx)
		if _, err := datastore.Put(ctx, key, &theClient.datastoreClient); err != nil {
			return Client{}, err
		}
		return theClient.datastoreClient.Client, nil
	}
	theClient.datastoreClient.fetched = true
	e, err := url.Parse(theClient.datastoreClient.Endpoint)
	if err != nil {
		return Client{}, err
	}
	client := New(theClient.datastoreClient.Key, theClient.datastoreClient.Secret).
		WithEndpoint(e).
		WithTransportFunc(newURLFetchTransportFunc(ctx))
	theClient.datastoreClient.Client = *client
	return theClient.datastoreClient.Client, nil
}

func placeholderDatastoreClient(ctx context.Context) datastoreClient {
	client := New("not-the-real-key", "not-the-real-secret").
		WithTransportFunc(newURLFetchTransportFunc(ctx))

	return datastoreClient{
		fetched:  true,
		Key:      client.key,
		Secret:   client.secret,
		Endpoint: client.endpoint.String(),
		Client:   *client,
	}
}
