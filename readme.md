# dean
server end for the vote app.

## build and run
bee run -gendoc=true -downdoc=true

## compile for production
GOOS=windows GOARCH=386 bee pack

## doc url
http://localhost:8090/swagger/

## sample url
http://localhost:8090/v1/dean/vote


## Design Considerations

## overall progress
refer at: