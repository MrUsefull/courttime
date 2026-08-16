# courttime

Court Time is a go package aimed around generating schedules for cross division play with limited court or field resources.

In other words, given

**Division**

- A
- B
- C
- D

**Constraints**

- Where A can play A or B teams
- B can play A, B, or C teams
- C can play B, C, or D teams
- D can play C or D teams

And there are court/field time constraints of
Court 1...N is available for M hours.
Each game takes T time (ie, there are some number of buckets to slot games into)

Provide a schedule for play, where each team gets "roughly" equal play time.

## Roadmap

- [ ] Basic multi court scheduling for cross division play
- [ ] Court affinity - ie, don't schedule back to back games for a team across two courts
- [ ] Team Rosters and subs (perhaps just left to do by hand?)
