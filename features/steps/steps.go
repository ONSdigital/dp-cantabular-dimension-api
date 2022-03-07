package steps

import (
	"fmt"
	"net/http"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/cucumber/godog"
)

func (c *Component) RegisterSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^the service starts`, c.theServiceStarts)
	ctx.Step(`^private endpoints are enabled`, c.privateEndpointsAreEnabled)
	ctx.Step(`^private endpoints are not enabled`, c.privateEndpointsAreNotEnabled)
	ctx.Step(`^cantabular server is healthy`, c.cantabularServerIsHealthy)
	ctx.Step(`^cantabular api extension is healthy`, c.cantabularAPIExtIsHealthy)
	ctx.Step(`^the following geography query response is available from Cantabular api extension for the dataset "([^"]*)":$`, c.theFollowingCantabularResponseIsAvailable)
	ctx.Step(`^the following error json response is returned from Cantabular api extension for the dataset "([^"]*)":$`, c.theFollowingCantabularResponseIsAvailable)
	ctx.Step(`^the following area query response is available from Cantabular api extension for the dataset "([^"]*)":$`, c.theFollowingCantabularAreaResponseIsAvailable)
	ctx.Step(`^the following area query response is available from Cantabular api extension for the dataset "([^"]*)" and text "([^"]*)":$`, c.theFollowingCantabularAreaTextResponseIsAvailable)

}

// theServiceStarts starts the service under test in a new go-routine
// note that this step should be called only after all dependencies have been setup,
// to prevent any race condition, specially during the first healthcheck iteration.
func (c *Component) theServiceStarts() error {
	c.startService(c.ctx)
	return nil
}

// cantabularServerIsHealthy generates a mocked healthy response for cantabular server
func (c *Component) cantabularServerIsHealthy() error {
	const res = `{"status": "OK"}`
	c.CantabularSrv.NewHandler().
		Get("/v9/datasets").
		Reply(http.StatusOK).
		BodyString(res)
	return nil
}

// cantabularAPIExtIsHealthy generates a mocked healthy response for cantabular server
func (c *Component) cantabularAPIExtIsHealthy() error {
	const res = `{"status": "OK"}`
	c.CantabularApiExt.NewHandler().
		Get("/graphql?query={}").
		Reply(http.StatusOK).
		BodyString(res)
	return nil
}

func (c *Component) privateEndpointsAreEnabled() error {
	c.cfg.EnablePrivateEndpoints = true
	return nil
}

func (c *Component) privateEndpointsAreNotEnabled() error {
	c.cfg.EnablePrivateEndpoints = false
	return nil
}

// theFollowingCantabularResponseIsAvailable generates a mocked response for Cantabular Server
// POST /graphql with the provided query and Cantabular dataset
func (c *Component) theFollowingCantabularResponseIsAvailable(dataset string, cb *godog.DocString) error {
	data := cantabular.QueryData{
		Dataset: dataset,
	}

	b, err := data.Encode(cantabular.QueryGeographyDimensions)
	if err != nil {
		return fmt.Errorf("failed to encode GraphQL query: %w", err)
	}

	// create graphql handler with expected query body
	c.CantabularApiExt.NewHandler().
		Post("/graphql").
		AssertBody(b.Bytes()).
		Reply(http.StatusOK).
		BodyString(cb.Content)

	return nil
}

func (c *Component) theFollowingCantabularAreaResponseIsAvailable(dataset string, cb *godog.DocString) error {
	data := cantabular.QueryData{
		Dataset: dataset,
	}

	b, err := data.Encode(cantabular.QueryAreasByArea)
	if err != nil {
		return fmt.Errorf("failed to encode GraphQL query: %w", err)
	}

	// create graphql handler with expected query body
	c.CantabularApiExt.NewHandler().
		Post("/graphql").
		AssertBody(b.Bytes()).
		Reply(http.StatusOK).
		BodyString(cb.Content)

	return nil
}

func (c *Component) theFollowingCantabularAreaTextResponseIsAvailable(dataset string, text string, cb *godog.DocString) error {
	data := cantabular.QueryData{
		Dataset: dataset,
		Text:    text,
	}

	b, err := data.Encode(cantabular.QueryAreasByArea)
	if err != nil {
		return fmt.Errorf("failed to encode GraphQL query: %w", err)
	}

	// create graphql handler with expected query body
	c.CantabularApiExt.NewHandler().
		Post("/graphql").
		AssertBody(b.Bytes()).
		Reply(http.StatusOK).
		BodyString(cb.Content)

	return nil
}
