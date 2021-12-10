#!/bin/bash

ab -n 100000 -c 200 -H "Accept:application/json" -k -p ./sale.json -T "application/json" http://localhost:8080/register/sale/