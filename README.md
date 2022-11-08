# Amazon Suppressed Check
Amazon Suppressed Check is a tool to add products and see if they are search suppressed or not. Not even that stupid arbitrage scripts could gave us if item is search suppressed, so here I am. Doing a useful shit, what a surprise.

Workflow:

1. Unmarshal 5 main league teams to a struct.
2. Get Amazon Client Credentials.
3. Declare an expiration timer to get client credentials once expire.
4. Declare logger.
5. Search product name from CSV in 5 structs to find a match, if match split the product name's prefix and suffix.
6. Request items data from /catalog/2022-04-01/items
7. If items data from endpoint contains either team prefix or suffix, get ASIN and product type of item.
8. Request restrictions from /listings/2021-08-01/restrictions, if no restrictions proceed.
9. Make a request to /listings/2021-08-01/items/%s/%s, where first parameter is seller ID, and second is SKU, marshal the body query and hash it.
10. Detect if there is any issues, by if issues are non-empty or summary is empty (not BUYABLE,DISCOVERABLE etc.) If it is problematic, then it is either incompleted, or search suppressed.

TODO: Declare tempASIN and productName from csv file, so automate it.

TODO: Change workflow depending on number of search results from /catalog/2022-04-01/items


