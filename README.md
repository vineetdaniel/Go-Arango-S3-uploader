A small utility script to recursively fetch images from given path and uploads them to S3. The meta information images is also stored in ArangoDB. The purpose of storing meta data to DB is to limit requests to S3 and to make searching for a particular images easier.

TO-DO 
 * Add initial run and incremental run options to avoid duplicate images and to make the script act like rsync.
 * Add UI and API for searching of images.
 * Add Analytics as to which images are being search/downloaded by adding hooks to S3.
 
 
