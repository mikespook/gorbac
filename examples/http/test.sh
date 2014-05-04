#!/usr/bin/env bash
# Add a role(r-a) with permissions [p-a, p-b]
curl -X PUT -d "name=r-a&permissions=p-a&permissions=p-b" http://127.0.0.1:12345/rbac
printf "\n"
# Add a role(r-b) with permissions [p-c, p-d]
curl -X PUT -d "name=r-b&permissions=p-c&permissions=p-d" http://127.0.0.1:12345/rbac
printf "\n"
# Add a role(r-c) with permission(p-e) and parents [r-a, r-b]
curl -X PUT -d "name=r-c&permissions=p-e&parents=r-a&parents=r-b" http://127.0.0.1:12345/rbac
printf "\n"
# test a permission(p-c) for a role(r-c)
curl -X GET "http://127.0.0.1:12345/isgranded?name=r-c&permission=p-c"
printf "\n"
# test a permission(p-x) for a role(r-c)
curl -X GET "http://127.0.0.1:12345/isgranded?name=r-c&permission=p-x"
printf "\n"
# remove a role(r-b)
curl -X DELETE "http://127.0.0.1:12345/rbac?name=r-b"
printf "\n"
# test a permission(p-c) for a role(r-c) again
curl -X GET "http://127.0.0.1:12345/isgranded?name=r-c&permission=p-c"
printf "\n"
# add a permission(p-c) to the role(r-a)
curl -X PATCH -d "name=r-a&permissions=p-c" http://127.0.0.1:12345/rbac
printf "\n"
# test a permission(p-c) for a role(r-c) again
curl -X GET "http://127.0.0.1:12345/isgranded?name=r-c&permission=p-c"
printf "\n"
# get all data
curl -X GET "http://127.0.0.1:12345/rbac"
printf "\n"
