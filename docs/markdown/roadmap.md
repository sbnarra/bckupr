# Roadmap

## Doing Now

* Complete Documentation
    * home badges... ci/cd - docker pulls/scans
* ~~scratch image, add 'mgr' component as entrypoint~~
    * ~~name, will house: cron, api, ui? mgr  web~~
* ~~image names bckupr/cli AND bckupr/cron AND bckupr/web ...~~
* Filesystem backup/restore (half implemented)

## Will Do Soon

* Backups need type prepend, full and partial (when filtering applied)
    * `YYYY-MM-DD_hh-mm-<type>` -full-cron -full-manual (ah not too sure tbh)
* support incremental backups

* opencontainer annotations added to image
    * https://github.com/opencontainers/image-spec/blob/main/annotations.md
* hide email? too much name in there


--- can publish at this point ---

* Metrics - 
    * after API to have long running process to serve metrics from?
    * prometheus - https://github.com/prometheus/client_python

* CI/CD - Github actions 
    * monitor base image, auto rebuild whenever there's a new base image

## Will Do Eventually

* API
    * POST /backups - trigger new backup
    * GET /backups - list backups
    * GET /backups/{id} - download backup
    * DELETE /backups/{id} - delete local backup


## Might Do Eventually

* Simple GUI - 
    * 1 page with list of backups
    * buttons to download/delete. 
        * toggles/selecting all
    * header button trigger new backups, 
    
## Could Do Someday

* ~~Rewrite the whole thing in GO? ... why? Go seems like a better lang for the use-case of writing docker apps/tools, need to research a bit more to decide~~
    * ~~API side of things looks good~~
    * ~~UI? guess it offers a web server so can stick to HTML/JS~~

~~GO will support running on the scratch image, python binary needs same image as builder image...might be better to do this now before writing any more code~~

~~https://medium.com/geekculture/how-to-structure-your-project-in-golang-the-backend-developers-guide-31be05c6fdd9~~

~~https://github.com/golang-standards/project-layout~~