# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: login.spec.ts >> login, post, comment, bookmark, and profile interactions
- Location: e2e\login.spec.ts:3:1

# Error details

```
Test timeout of 30000ms exceeded.
```

# Page snapshot

```yaml
- generic [ref=e1]:
  - generic [ref=e3]:
    - banner [ref=e4]:
      - generic [ref=e5]:
        - generic [ref=e6]: CloudCode SNS
        - button "Logout" [ref=e7] [cursor=pointer]
    - list [ref=e12]:
      - listitem [ref=e13]:
        - button "Timeline" [ref=e14] [cursor=pointer]:
          - img [ref=e16]
          - generic [ref=e19]: Timeline
      - listitem [ref=e20]:
        - button "Bookmarks" [ref=e21] [cursor=pointer]:
          - img [ref=e23]
          - generic [ref=e26]: Bookmarks
      - listitem [ref=e27]:
        - button "Profile" [active] [ref=e28] [cursor=pointer]:
          - img [ref=e30]
          - generic [ref=e33]: Profile
    - main [ref=e34]:
      - generic [ref=e37]:
        - generic [ref=e39]:
          - generic [ref=e42]:
            - generic [ref=e43]: OR
            - generic [ref=e45]:
              - heading "Orange Peach" [level=4] [ref=e46]
              - paragraph [ref=e47]: orange.peach.695@example.com
              - generic [ref=e48]:
                - generic [ref=e49] [cursor=pointer]: 0Followers
                - generic [ref=e50] [cursor=pointer]: 0Following
          - tablist "profile tabs" [ref=e54]:
            - tab "Posts" [selected] [ref=e55] [cursor=pointer]
            - tab "Followers" [ref=e56] [cursor=pointer]
            - tab "Following" [ref=e57] [cursor=pointer]
          - generic [ref=e60]:
            - generic [ref=e61]:
              - generic [ref=e63] [cursor=pointer]: OR
              - generic [ref=e64]:
                - heading "Orange Peach" [level=6] [ref=e65] [cursor=pointer]
                - text: 5/27/2026, 11:04:45 PM
              - button [ref=e67] [cursor=pointer]:
                - img [ref=e68]
            - paragraph [ref=e71]: "E2E testing new premium features: follows and details!"
            - generic [ref=e72]:
              - button "like" [ref=e73] [cursor=pointer]:
                - img [ref=e75]
              - button "comment" [ref=e77] [cursor=pointer]:
                - img [ref=e78]
              - button "bookmark" [ref=e80] [cursor=pointer]:
                - generic [ref=e81]:
                  - img [ref=e82]
                  - generic [ref=e84]: "1"
              - button "share" [ref=e85] [cursor=pointer]:
                - img [ref=e86]
        - generic [ref=e89]:
          - generic [ref=e90]:
            - heading "Trending Topics" [level=6] [ref=e91]
            - list [ref=e92]:
              - listitem [ref=e93]:
                - button "#CloudCode 10.5k posts" [ref=e94] [cursor=pointer]:
                  - generic [ref=e95]:
                    - generic [ref=e96]: "#CloudCode"
                    - paragraph [ref=e97]: 10.5k posts
              - listitem [ref=e98]:
                - button "#GoLang 8.2k posts" [ref=e99] [cursor=pointer]:
                  - generic [ref=e100]:
                    - generic [ref=e101]: "#GoLang"
                    - paragraph [ref=e102]: 8.2k posts
              - listitem [ref=e103]:
                - button "#ReactMUI 3.1k posts" [ref=e104] [cursor=pointer]:
                  - generic [ref=e105]:
                    - generic [ref=e106]: "#ReactMUI"
                    - paragraph [ref=e107]: 3.1k posts
          - generic [ref=e108]:
            - heading "Who to follow" [level=6] [ref=e109]
            - paragraph [ref=e110]: Suggestions will appear here based on your activity.
            - button "Find People" [ref=e111] [cursor=pointer]
  - iframe [ref=e112]:
    
```