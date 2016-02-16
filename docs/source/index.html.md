---
title: API Reference

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
The API server manages a munki repo over HTTP endpoints, updating the underlying datastore, which can be one of:  

* a plain collection of plist files on disk  
* plist files on disk, version controlled in a git repo  
* a database  

## Status

### Endpoints
/manifests - documented below
/pkgsinfo - should work the same as manifests
/pkgs/*path/to/pkg - POST or DELETE binary files

### Content and Accept Headers

All requests support `application/xml`(plist) and `application/json` as `Accept` and `Content-Type` headers. 


> To send a plist

```shell
curl -H "Content-Type: apllication/xml" -d@myplist.plist http://api/endpoint
```

> To receive response as a plist:  

```shell
curl -H "Accept: application/xml" -d@myplist.plist http://api/endpoint
```
# Authentication

None Yet. You can implement authentication by using nginx or apache as a proxy.
Check out [Kong](https://getkong.org/plugins/)


# Manifests

## Get All Manifests

```shell
curl "http://example.com/api/manifests"
```
> The above command returns JSON structured like this:

```json
[
 {
  "filename": "C02KP0H4DR55",
  "catalogs": [
   "production"
  ],
  "display_name": "John's MBP",
  "included_manifests": [
   "admin-common"
  ],
  "managed_installs": [
   "PuppetAgent"
  ],
  "notes": "Some Admin Notes",
  "user": "John Doe"
 },
 {
  "filename": "vagrant-15C50",
  "catalogs": [
   "testing"
  ],
  "optional_installs": [
   "AdobePhotoshop"
  ],
  "managed_installs": [
   "GoogleChrome",
   "munkitools",
   "munkitools_admin",
   "munkitools_core",
   "munkitools_launchd",
   "sal_scripts"
  ]
 }
]
```

This endpoint retrieves all the manifests in a repo.

### HTTP Request

`GET http://example.com/api/manifests`

## Get a Specific Manifest

```shell
curl "http://example.com/api/manifests/site_default"
```

> The above command returns JSON structured like this:

```json
{
 "catalogs": [
  "production"
 ],
 "managed_installs": [
  "munkitools",
  "munkitools_core",
  "munkitools_launchd"
 ]
}
```

This endpoint retrieves a specific manifest.

### HTTP Request

`GET http://example.com/manifests/:name`

## Create a new manifest

```shell
curl -H "Content-Type: application/json" \
     -X POST --data \
`{
 "filename": "foo.example.com",
 "catalogs": [
  "production"
 ],
 "display_name": "example manifest"
}`
```

> On success, the above command returns HTTP 201 Created and
> JSON structured like this:

```json
{
 "catalogs": [
  "production"
 ],
 "display_name": "example manifest"
}
```

> If the resource already exists, the API will return 409 Conflict

This endpoint creates a new manifest.

### HTTP Request

`POST http://example.com/manifests`

## Delete a manifest

```shell
curl -X DELETE "http://example.com/api/manifests/foo.example.com"
```

> On success, the above command returns HTTP 204 No Content

This endpoint deletes an existing manifest.

### HTTP Request

`DELETE http://example.com/manifests/:name`

## Update an existing manifest

```shell
curl -H "Content-Type: application/json" \
     -X PATCH --data \
`{
 "display_name": "updated manifest"
}` http://example.com/manifests/foo.example.com
```

> On success, the above command returns HTTP 202 Ok and
> JSON structured like this:

```json
{
 "catalogs": [
  "production"
 ],
 "display_name": "updated manifest"
}
```

This endpoint updates an existing manifest

### HTTP Request

`PATCH http://example.com/manifests/:name`

# PkgsInfos

works like manifests, needs to be documented

# Pkgs
Binary files can be added and deleted:

```shell
curl -vX POST http://example.com/api/pkgs/apps/Firefox-43-0-4.dmg --data-binary @Firefox-43.0.4.dmg
```

```shell
curl -vX DELETE http://example.com/api/pkgs/apps/Firefox-43-0-4.dmg
```
