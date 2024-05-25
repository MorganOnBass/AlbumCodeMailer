>Album Code Mailer

This is a small utility for using AWS SES to send album codes to patrons and/or crowdfunders in templated emails.

You will need to ensure your AWS account meets the requirements to make the SendBulkTemplatedEmail api call: https://docs.aws.amazon.com/ses/latest/dg/send-personalized-email-api.html

Templating is fairly straightforward, compose your email and use the form `{{code}}` where you would like it to appear.
Example:

```
Dear patron,

Thanks for the beer money, go get your dank tunes with the code {{code}} 
at totallyrealbandname.bandcamp.com/whatever...
```

The UploadTemplate command will want a path to an email template in both plaintext and html form. 

The recipients input file expects something like what I've observed kickstarter and patreon exports to look like - 
a CSV with column headers on the first row and email addresses in the second column. Like:

```
Name,Email
Alice,example@example.com
Bob,example2@example2.com
```

And the codes list should be a plain text file with one code per line, like: 

```
0123-4567
abcd-ef01
...
```

Once this is all set up, upload the template like so:

```
go run main.go UploadTemplate -n AwesomeTemplateName \
-s 'Here is your download code for some dank tunes!' \
-t /Users/morgan/code/AlbumCodeMailer/Templates/Patreon/textpart.txt \
-p /Users/morgan/code/AlbumCodeMailer/Templates/Patreon/htmlpart.html
```

Then you should test your template with a single address:

```
go run main.go TestBulkTemplatedEmail -d destination.address@domain.com \
-f from.address@domain.com -s replyto.address@domain.com -t AwesomeTemplateName
```

When you're ready to send the lot, do a dry run to sanity check the output. A dry run prints each address and the code 
that will be sent to it on one line. Ensure this looks right - valid codes, unique codes, etc.

```
go run main.go SendBulkTemplatedEmail -c /Path/To/List_of_codes \
-i /path/to/recipients.csv -f from.address@domain.com -s reply-to.address@domain.com \
-t AwesomeTemplateName -d
```

If that checks out, simply run the above again without the trailing `-d` to begin sending emails. To avoid rate limiting,
this program sends 100 emails per minute. This is very conservative, but I hacked this together in an evening and I
don't feel like doing that math. If you have enough patrons or backers for this to hurt, I guess look to optimise here.
Or just get your publicist to handle it. ;)