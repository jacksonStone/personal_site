# How I Host This, and All My Hobby Projects on the Cheap

## The Problem

I am, at my professional core, a full stack developer. I love the satisfaction that comes from starting with a blank text file and ending with an end-to-end experience. This means if I do a hobby project, it's generally a website of some kind hosted on the World Wide Web. However, if you are anything like me, your hobby projects suffer from a similar fate:

- Hobby project costs X dollars a month to keep running
- You work on them in bursts - sometimes spaced out by years
- Every month when the AWS/Heroku/whatever bill rolls in, it's time to think: "Do I really want to pay this fee to keep this toy thing running?"

This has led to me canceling many previous sites, which subsequently leads to letting a domain name expire. (Which I always regret eventually)

This has long plagued me. In addition to leading to an early fate for my hard work, it keeps me from starting otherwise fun and fruitful projects that would have been nice learning experiences and tiny contributors to my overall portfolio. The thought pattern works something like this:

> "I will cancel it in a month and will have nothing to show but Github links for it. Might as well not start."

It's not that I don't enjoy programming for its own sake, but the foregone conclusion I won't have a work to show at the end is sufficiently discouraging to prevent me from getting started. An analogy might be if you were a prolific wood craftsperson but knew at the end you were going to chuck it into the fire. Part of the joy is being able to look back at your creation even if you love the craft for its own sake.

Additionally, if a wood craftsperson were in a position where they would have to pay 10-20 bucks a month to let the craft live in their house, they too would be tempted every so often to get rid of it. If we are developing sites for fun, that has this unfortunate effect. It isn't even that the money sums are large or consequential. It's simply knowing there is a slow drip and a constant cost/benefit analysis being made.

I had to come up with a new way. I wanted to keep up the projects, but not feel pressured to throw them away. And if I could avoid that spiral, I think it would lead to more hacking tendencies that I think have been key to my career success thus far. (Keeping skills sharp, passion alive for the craft, etc.)

## Tech

The general solution I've discovered for this problem is the following:

1. Host everything on one small server with a static IP address. (In my specific case I reserved 3 years of a t3.small EC2 for about a 50% discount. It averages about $7 a month amortized out.) This server is more than enough for hobby projects which get visitors in the low double digits monthly if I am honest with myself. :D
2. Point all domains for all projects to this one static IP address.
3. Use a small reverse proxy server running on the box ([around 100 lines of Go](https://github.com/jacksonStone/little_reverse_proxy/blob/main/server.go)) that handles HTTPS resolution and redirects traffic to a local host port based on the domain name. (Can't have my hobby projects on HTTP like a savage! If you show up at one of my sites - you're getting your data encrypted!) It zips along and uses basically no system resources as opposed to some more fully featured alternatives.
4. Configure autocert on the box for each domain to handle SSL Cert rotation, which also restarts the proxy when there is a new cert. (HTTPS will then "Just work" and I won't have a broken site once every year or so... I want these hobby projects to last decades.)
5. Spin up each hobby project on its own localhost port on this box (which the reverse proxy knows to point to).
6. In each side project, build a deploy.sh script (like the one for [this site](https://github.com/jacksonStone/personal_site/blob/main/deploy.sh)) that with one command - will handle transferring files to the server (before or after building locally depending on the tech stack) and restart the service.
7. Create simple .service files on the ubuntu box to run the executables on startup/reboot/whatever so even if something happens in the Data center holding my service and AWS needs to rotate it - I'll still be sitting pretty.

> *As an aside: You know who says treat servers like cattle? Those who sell you servers...*

## Cost

In total my hosting costs for all of this is broken down as follows:

1. EC2 T3.Small running Ubuntu: about $7/month with a 3-year reservation
2. One "Elastic IP" assigned to this box (An IPV4 address): about $5/month (crazy that it is as costly as the server, but there you go...)
3. Free tier MongoDB for the NoSQL apps ($0... for now at least)
4. SQLite instances for SQL related side projects on the ubuntu box (0.08GB/Month of EBS - Essentially $0/month)
5. The domain names renewal costs (varies depending on the Top Level Domain (.com/.info/.whatever))
   - Maybe anywhere from $3-5/month - paid about once a year
   - Plus $0.50/month per hosted zone in AWS

With all this put together, the majority of the expense is the once-a-year cost of the domain names, not the actual server, etc. Each new side project in effect only adds the cost of its own domain name. There is essentially no incremental cost beyond that (so long as the project is not RAM hungry or something justifying a bigger box).

## Performance

Given the presumption we are dealing with hobby projects, the performance considerations do not really need to account for scaling beyond, at the high end, GBs. Considering I am not using a CDN, any visitor to my sites in mainland US is looking at about a 40ms round trip time. This ballpark is a combination of longest path between west coast and east coast as well as assuming optimal network cabling between those two destinations, given the speed of light.  This is acceptable for landing pages, etc. It assumes of course I don't do anything in my services to poopoo the performance, as everything will pay that 40ms cost. 

Once on the box, the reverse proxy itself seems to take on the order of 0.5-1ms. (Given the amount of time it takes to reroute the request to the correct port on the box). I even do a small optimization of going [direct to 127.0.0.1 to avoid DNS resolution.](https://www.youtube.com/watch?v=Pfy4Q-uDV6I) Though the difference is not easily measured. 

SQLite writes recrods and reads them on the order of 1-10ms since it never leaves the box. And MongoDB free tier seems to return a response anywhere from 100-500ms. (This is enough to frustrate me so I will likely migrate everything to SQLite in the future... I have also learned to distrust all "Free" tier services. Looking at you Heroku!) 

After this, it boils down to whatever other computation I am doing in the services to massage or otherwise alter the payload. This is generally 1-10ms for most of my endpoints.

In Summary, the majority of performance considerations for the sites are the time it takes physically for light to travel to my box and back. And if it uses MongoDB, to then have it go to Mongo and back. I can live with that. CDN, shme-DN.

## Conclusion

It's laughable how much time I spent getting this cost down to as low as I did. (Maybe 4-5 days, though to be fair it was while on Baby leave in between diaper changes for my third baby) But as mentioned earlier - it's not really about the money, but the constant temptation to turn down the hobby sites due to ongoing expense. I wanted to minimize incremental cost for a project and essentially become my own tiny cloud provider ensuring my hobby projects could have a long and happy life on the shelf.

At present of writing this post, I am hosting four services on this one box I've configured (NodeJS full stack apps and Go apps) and I'm hovering at about 20% memory utilization in RAM and CPU utilization percentage is low single digits. So plenty of headroom for something like 10 more services. Additionally - it's easier for me to stomach the cost of hosting a single server because the question is now "Do I want to host anything?" rather than little hobby site X or Y. This is more similar to a wood craftsperson paying for their tools every couple years, rather than paying for each piece of woodworking they have in their house. Easy choice.

Anyway, maybe if you are like me you will find my musings useful. Let's go hack together some projects!