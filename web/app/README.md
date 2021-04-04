# ls

## Build Setup

```bash
# install dependencies
$ yarn install

# serve with hot reload at localhost:3000
$ yarn dev

# build for production and launch server
$ yarn build
$ yarn start

# generate static project
$ yarn generate
```

For detailed explanation on how things work, check out [Nuxt.js docs](https://nuxtjs.org).


<form id="upload" method="post"
      action="/couchapp/doc_with_attachment"
      enctype="multipart/form-data">
    <label>Revision : <input id="revision" type="text" name="_rev"/></label><br/>
    <input id="attachment" type="file" name="_attachments"/><br/>
    <br/>
    <input type="submit"/>
</form>