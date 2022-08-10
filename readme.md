# Neo4j - LinkedIn
After listening many times about **graph databases**, I was thinking in possibles scenarios where I could take advantage
of them. Then I came up with this idea that could be useful in two possibles ways:
- Keep track of your **profile viewers**, would allow you to store who visited your profile and the last time they did it. 
- If used globally, with a common database, the people could register their interactions and this way you could see if 
 a certain recruiter is really interested with you or if it's speaking with everybody. So you could use it as
**spam filter**.

![Example graph](/images/example-graph.svg)

For the graph shown above, if we request Sam's profile viewers, we would get the following response: 

```json
[
    {
        "Person": {
            "Name": "Olivia",
            "Skills": null
        },
        "When": "2022-08-10T18:35:42.937Z"
    },
    {
        "Person": {
            "Name": "Tom",
            "Skills": null
        },
        "When": "2022-08-10T17:36:41.052Z"
    }
]
```

Has been a great experience and the Cypher Query Language it's simple and very cool!
```cypher
MATCH (n)-[r:VIEWED_PROFILE]->(p:Person {name: $name})
RETURN n.name, n.skills, r.datetime
```
