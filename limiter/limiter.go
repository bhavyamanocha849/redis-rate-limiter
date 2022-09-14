/*
*Its the decision engine
*Get the request
*Fetch the configuration from the cache/nosql db
*Pass the values to the decision engine to check if the request is valid or not
*if valid{
*	update the cache
TODO:Check the type of caches
*}else{
*	Reject the request with 429 error
*}
*/
package limiter

import (
	"context"

	"github.com/bhavyamanocha849/redis-rate-limiter/datatypes"
)

type Limiter interface {
	Run(ctx context.Context, response *datatypes.Request) (*datatypes.Response, error)
}
