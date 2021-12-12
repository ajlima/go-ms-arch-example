#!/bin/bash
ab -n 1000000 -c 250 -H "Accept:application/json" -k -p ./sale.json -T "application/json" http://localhost:8081/api/v1/register/sale/