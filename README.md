EyeGo
=====

Eyefi Upload Server written in Go. 

* Full support of photo & movie uploading
* Processes files then moves them to target directory
* Geotagging support via the Google Geolocation API (max. 100 requests / day)
* BSD licensed

Requires a configuration file in JSON, using the following format:

```
{
    "target_dir":"/Users/jpg/Photos",
    "google_api_key":"XXX",
    "cards":[
        {
            "mac_address":"XXX",
            "upload_key":"XXX"
        }
    ]
}
```
