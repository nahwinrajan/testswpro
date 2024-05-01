Hi Hilda and Team,

I've completed the task and pushed it to [GitHub](https://github.com/nahwinrajan/testswpro/tree/main). Please use the `main` branch for the submission.

**Assumption:**
The calculation in the problem statement is wrong because it does not count movement/distance traveled to the last plot. Therefore, my calculation counts it. Instead of 54 as in the problem statement, it is 64.

I prepared the database and logic to accommodate `max_distance`, the bonus question, but I did not present it (disabled the URL) because I don't have time to cover all the edge cases. The logic and database design in the answer do cater for such things, even the base code is in the code but commented.

There are edge cases where I am not sure what is expected, and even if I do, I no longer have time to finish it due to fixing mundane stuff from echo (e.g., binding to interface instead of instance, problematic usage of UUID type in PostgreSQL).

I wrote the unit test using vanilla golang style with `_test.go` and mocked interface instead of utilizing the `test/*` folder presented in SDK

Ideally, I would create an organization following clean code, separating handler, business logic, and data, but I'm not really sure on this one because I must work on SDK. I did most of the unit tests, including the repo layers, except for one handler `PostEstateIdTree`. I did write unit test for it, but it still lacks handling of couple last steps; thus I commented the unit test. All other functions that are used there are unit tested. Somehow go tool for testing is including the generated interface for coverage calculation. As stated before, there is only one function that is not having a functioning unit test, albeit all functions within those flows are tested. There are a couple `NEW` functions (server and repository) which do not really have anything to test. That being said, I believe the coverage must be at least 70% as requested.

If I were to organize it as I would, it will be:

- `cmd/web/main.go`: main package for building web
- `internal/handler`: for transport protocol layer, will convert payload into model understood throughout the project and back to respective response body required
- `internal/usecase`: for business logic
- `internal/repo`: for interaction with data storage
- `internal/model`: the DTO/structure for the data type that is needed

**FLOW:** `handler -> usecase -> repo`. Only from that way model is accessible throughout the internal folders

`pkg/*`: helper package needed goes here

Thank you for the chances, enjoyed jumping into echo framework instead of writing everything from scratch.

Looking forward to further discussion,

Cheers