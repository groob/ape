---
title: Munki API Reference

language_tabs:
  - shell

toc_footers:
  - <a href='#'>Sign Up for a Developer Key</a>
  - <a href='https://github.com/tripit/slate'>Documentation Powered by Slate</a>

includes:
  - errors

search: true
---

# Introduction

Welcome to Ape, an API for Munki.
The API server manages a Munki repo over HTTP endpoints, updating the underlying datastore.

# Authentication

None yet. You can implement authentication by using nginx or apache as a proxy.
Check out [Kong](https://getkong.org/plugins/)


# Manifests

## Get all manifests

```shell
curl -X GET \
     -H "Content-Type: application/json" \
     -H "Accept: application/json" \
     http://example.com/api/manifests
```

> With parameters: 

```shell
curl -X GET \
     -H "Content-Type: application/json" \
     -H "Accept: application/json" \
     http://example.com/api/manifests?catalogs=production&api_fields=filename,managed_installs
```

This endpoint retrieves all the manifests in a repo.

### HTTP Request

`GET http://example.com/api/manifests`

### Query Parameters
Parameters can be used to filter the results by one or more manifest keys.  
The <code>api_fields</code> key can be used to restrict which fields should be returned by the API.

**Example:**

Parameter |  Description
--------- | -----------
catalogs=testing,dev| Filter by catalog, using commas to separate catalogs
api_fileds=filename,managed_installs | Only return the <code>filename</code> and <code>managed_installs</code> fields.

## Get a specific manifest

```shell
curl -X GET \
     -H "Content-Type: application/json" \
     -H "Accept: application/json" \
     http://example.com/api/manifests/site_default
```

This endpoint retrieves a specific manifest.

### HTTP Request

`GET http://example.com/api/manifests/:filename`

## Create a new manifest

```shell
curl -X POST \
     -H "Content-Type: application/json" \
     -H "Accept: application/json" \
     -d @Firefox-43.0.4.json http://example.com/api/manifests
```

> On success, the above command returns HTTP 201 Created and the created manifest in the HTTP response body

> If the resource already exists, the API will return 409 Conflict

<aside class="notice">
A POST request to create a manifest must have a <code>filename</code> key in the request body.<br/>
The filename key should be the filename of the manifest in the repository. <br />
Example: <br />
<code>{ "filename" : "site_default", }
</code>
</aside>

This endpoint creates a new manifest.

### HTTP Request

`POST http://example.com/api/manifests`

## Delete a manifest

```shell
curl -X DELETE "http://example.com/api/manifests/site_default
```
> On success, the above command returns HTTP 204 No Content

This endpoint deletes an existing pkgsinfo.

### HTTP Request

`DELETE http://example.com/api/manifests/:filename`

## Update an existing manifest

```shell
curl -X PATCH \
     -H "Content-Type: application/json" \
     -H "Accept: application/json" \
     -d site_default.json http://example.com/api/manifests/site_default
```

> On success, the above command returns HTTP 202 Ok and the updated manifest in the HTTP response body.

This endpoint updates an existing manifest and can be used to send data for only the specific fields that you intend to update.

### HTTP Request

`PATCH http://example.com/api/manifests/:filename`

## Replace an existing manifest

```shell
curl -X PUT \
     -H "Content-Type: application/json" \
     -H "Accept: application/json" \
     -d site_default.json http://example.com/api/manifests/site_default
```

> On success, the above command returns HTTP 202 OK and the updated manifest in the HTTP response body.

This endpoint replaces an existing manifest with a new one.

### HTTP Request

`PATCH http://example.com/api/manifests/:filename`

# Pkgsinfos

## Get all pkgsinfos

```shell
curl -X GET \
     -H "Content-Type: application/json" \
     -H "Accept: application/json" \
     http://example.com/api/pkgsinfo
```

This endpoint retrieves all the pkgsinfos in a repo.


### HTTP Request

`GET http://example.com/api/pkgsinfo`

### Query Parameters
Parameters can be used to filter the results by one or more pkgsinfo keys.  
The <code>api_fields</code> key can be used to restrict which fields should be returned by the API.

**Example:**

Parameter |  Description
--------- | -----------
name=a,b,c| Filter by name, using commas to separate the pkg names
catalogs=testing,dev| Filter by catalog, using commas to separate catalogs
api_fileds=filename,name | Only return the <code>filename</code> and <code>name</code> fields.

## Get a specific pkgsinfo

```shell
curl -X GET \
     -H "Content-Type: application/json" \
     -H "Accept: application/json" \
     http://example.com/api/pkgsinfo/apps/Firefox-43.0.4.plist
```

This endpoint retrieves a specific pkgsinfo file.

### HTTP Request

`GET http://example.com/api/pkgsinfo/*filepath`

## Create a new pkgsinfo

```shell
curl -X POST \
     -H "Content-Type: application/json" \
     -H "Accept: application/json" \
     -d @Firefox-43.0.4.json http://example.com/api/pkgsinfo
```

> On success, the above command returns HTTP 201 Created and the created pkgsinfo in the HTTP response body

> If the resource already exists, the API will return 409 Conflict

<aside class="notice">
A POST request to create a pkgsinfo must have a <code>filename</code> key in the request body.<br/>
The filename key should be the relative path of the pkgsinfo in the repository. <br />
Example: <br />
<code>{ "filename" : "apps/Firefox-43.0.4.plist", }
</code>
</aside>

This endpoint creates a new pkgsinfo.

### HTTP Request

`POST http://example.com/api/pkgsinfo`

## Delete a pkgsinfo

```shell
curl -X DELETE "http://example.com/api/pkgsinfo/apps/Firefox-43.0.4.plist"
```

> On success, the above command returns HTTP 204 No Content

This endpoint deletes an existing pkgsinfo.

### HTTP Request

`DELETE http://example.com/api/pkgsinfo/*filepath`

## Update an existing pkgsinfo file

```shell
curl -X PATCH \
     -H "Content-Type: application/json" \
     -H "Accept: application/json" \
     -d @Firefox-43.0.4.json http://example.com/api/pkgsinfo/apps/Firefox-43.0.4.plist
```

> On success, the above command returns HTTP 202 Ok and the updated pkgsinfo in the HTTP response body

This endpoint updates an existing pkgsinfo and can be used to send data for only the specific fields that you intend to update.

### HTTP Request

`PATCH http://example.com/api/pkgsinfo/*filepath`

## Replace an existing pkgsinfo file

```shell
curl -X PUT \
     -H "Content-Type: application/json" \
     -H "Accept: application/json" \
     -d @Firefox-43.0.4.json http://example.com/api/pkgsinfo/apps/Firefox-43.0.4.plist
```

> On success, the above command returns HTTP 202 Ok and the updated pkgsinfo in the HTTP response body

This endpoint replaces an existing pkgsinfo with a new one

### HTTP Request

`PATCH http://example.com/api/pkgsinfo/*filepath`

# Pkgs

## Add a new pkg

```shell
curl -X POST \
     -F filename=apps/Firefox-43-0-4.dmg \
     -F filedata=@/path/to/local_file.dmg \
     http://example.com/api/pkgs
```

## Delete a pkg

```shell
curl -X DELETE http://example.com/api/pkgs/apps/Firefox-43-0-4.dmg
```

# Icons

## Add a new icon

```shell
curl -X POST \
     -F filename=Firefox.png \
     -F filedata=@/path/to/local_file.png \
     http://example.com/api/icons
```

## Delete an icon

```shell
curl -X DELETE http://example.com/api/icons/Firefox.png
```

# Media Types
The Munki API uses JSON by default but can also accept and return XML PLIST if specified in the request header.

Here are the headers you can pass to the API using `Content-Type` and `Accept` headers:  
`application/json`  
`application/json; charset=utf-8`  
`application/xml`  
`application/xml; charset=utf-8`  

# Error Response

In case of an error, the API will return one of the HTTP status codes below, as well as an `errors` object containing a list of errors in the response body.  
Example 404 response:  

```json
{
 "errors": [
  "Failed to fetch pkgsinfo from the datastore: Resource not found"
 ]
}
```

