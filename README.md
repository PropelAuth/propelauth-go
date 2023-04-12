todo:

- write the ValidateAccessTokenAndGet* methods
- write the (related) jwt library
- add tests

todo refactoring:

- rearrange the files into "models" "helpers" "jwt"
- organize endpoint methods into "endpoints" dir?
- test .github workflow
- some of the endpoints take in a couple parameters, some take in a params struct, make sure when we do one or the other is consistent. specifically around FetchBatchUserMetadataBy*.
- instal and run lint and fix recomendations
