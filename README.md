# `openApiCodeGen` Generate minimal javascript from a swagger file
Made for swagger v2 and outputs javascript that just calles a helper function for fetching the data 

## How to use

### Install
```sh
go get github.com/mjarkk/openApiCodeGen
```

### Generate javascript
```sh 
# Generate javascript and output it to the stdout
openApiCodeGen -in swagger.json

# Generate javascript and place it in a file
openApiCodeGen -in swagger.json > api.js
```

### Create the fetch function 
The api file generated will reqeust a function `apiUtil.js`  
The `apiUtil.js` file needs to be in the same dir as the generated code  
```js
// apiUtil.js
// apiFetcher is the function requested by the generated code

export const apiFetcher = async reqData => {
  console.log(reqData)
  /*{
    params: ['firstParam', 'secondParam'],
    method: "POST",
    url: "/api/v1/user/firstParam/edit/secondParam",
    body: {...},
  }*/

  return await fetch(reqData.url)
}
```
