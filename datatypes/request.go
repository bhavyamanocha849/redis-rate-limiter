package datatypes

/*
*Request is the actual request that will come
*We will only get the key
*The key will be in accordance with the api that we are limiting
 */

type Request struct {
	Key string //api key
}
