# Examples

Real-world examples of managing Reddit Ads campaigns from the terminal with `rad`.

---

## Browse your campaigns

```
$ rad campaign list
id                   name                       configured_status  effective_status  objective
-------------------  -------------------------  -----------------  ----------------  -----------
1234567890           Example Click Campaign     ACTIVE             ACTIVE            CLICKS
1234567891           Default Campaign           PAUSED             CAMPAIGN_PAUSED   IMPRESSIONS
```

Get the details of a campaign by name or ID:

```
$ rad campaign get "Example Click Campaign"
Campaign Example Click Campaign (1234567890)

{
  "configured_status": "ACTIVE",
  "effective_status": "ACTIVE",
  "objective": "CLICKS",
  "funding_instrument_id": "76401",
  "is_campaign_budget_optimization": false,
  "created_at": "2026-01-01T01:01:01.000000+00:00",
  ...
}
```

## Check ad group targeting

```
$ rad adgroup list
id                   name                         campaign_id          configured_status  effective_status
-------------------  ---------------------------  -------------------  -----------------  ----------------
1234567890           Example Traffic Campaign     123456789            ACTIVE             ACTIVE
```

## See how your ads are performing

Pull a campaign-level summary for the last 30 days:

```
$ rad report campaign-summary --since 30d
campaign_name              campaign_id          impressions  clicks  ctr    spend  cpc   ecpm
-------------------------  -------------------  -----------  ------  -----  -----  ----  ----
Example Click Campaign     123456789            812          13      0.016  6.66   0.51  8.20
```

Break it down by individual ad to find your winners:

```
$ rad report ad-summary --since 30d --campaign "Example Click Campaign"
ad_name                         impressions  clicks  ctr    spend  cpc   ecpm
------------------------------  -----------  ------  -----  -----  ----  ----
Traffic Feb 03, 2025 AAC ad 19  747          13      0.017  6.66   0.51  8.91
Traffic Feb 03, 2025 AAC ad 15  16           0       0      0.00   0.00  0.00
Traffic Feb 03, 2025 AAC ad 11  9            0       0      0.00   0.00  0.00
...
```

## Inspect ad creatives

See the headline, image, and CTA behind any ad:

```
$ rad post get t3_1ih2iit
{
  "headline": "Example ad headline",
  "type": "IMAGE",
  "content": [
    {
      "call_to_action": "View More",
      "destination_url": "https://www.example.com/?utm_source=reddit",
      "media_url": "https://reddit-image.s3.us-east-1.amazonaws.com/image.jpg"
    }
  ],
  "allow_comments": false,
  ...
}
```

## Create new ad creatives

Create a post with an image and headline:

```
$ rad post create \
    --profile abcdefg \
    --type IMAGE \
    --headline "My Headline" \
    --content-json '[{
      "call_to_action": "View More",
      "destination_url": "https://www.example.com/?utm_source=reddit",
      "media_url": "https://reddit-image.s3.us-east-1.amazonaws.com/image.jpg"
    }]' \
    --allow-comments false
```

## Attach ads to an ad group

Create an ad from a post and attach it to an ad group:

```
$ rad ad create \
    --ad-group "Example Traffic Campaign" \
    --name "Ad name" \
    --configured-status ACTIVE \
    --post-id abcdefg
Ad created.

{
  "id": "987654321",
  "name": "Ad name",
  "configured_status": "ACTIVE",
  "effective_status": "PENDING_APPROVAL",
  "ad_group_id": "123456789",
  "post_id": "abcdefg",
  ...
}
```

## Pause underperforming ads

```
$ rad ad update --configured-status PAUSED "ad name"
Ad updated: ad name (987654321)
```

## Update a campaign

Rename a campaign:

```
$ rad campaign update --name "Default Campaign" "Example Click Campaign"
Campaign updated: Example Click Campaign (123456789)
```

## Find communities to target

```
$ rad targeting communities search --query "3d printing"
name                  id         subscriber_count  categories    description
--------------------  ---------  ----------------  ------------  -------------------------------------------------
3Dprinting            t5_2rk5q   3315756           3D Printing   From models to figurines, explore the possibilit...
functionalprint       t5_30567   591181            3D Printing   Find inspiration for 3D-printed parts that are b...
resinprinting         t5_2ysf8   122461            3D Printing   Dive into the world of resin printing and bring ...
FixMyPrint            t5_30mvb   198617            3D Printing   Get expert help to troubleshoot and fix your 3D ...
BambuLab              t5_69mkea  391250            3D Printing   Join the conversation on BambuLab 3D printers...
PrintedMinis          t5_3i0p1   152580            Tabletop...   Find inspiration, resources, and advice on creat...
...
```

## Check your funding

```
$ rad funding list
name        id     currency  credit_limit  billable_amount  is_servable
----------  -----  --------  ------------  ---------------  -----------
Self-Serve  12345  USD       2500.00       6.66             true
```

## Export data

Any command supports `--json` for machine-readable output, and reports support `--csv`:

```
$ rad report campaign-summary --since 7d --csv > weekly-report.csv
$ rad campaign list --json | jq '.[].name'
"Example Click Campaign"
"Default Campaign"
```
