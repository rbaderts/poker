 

This is a work-in-progress multi-user holdem poker web application,
built primarily to experiment with new techologies and 
architecture patterns.

The backend is in golang, the front-end uses vue.js.     
The application employs the ideas from "Clean Archtiecture".   
It employ async push notifications to attached web clients over
websockets...


It uses github.com/rbaderts/pokerlib which provides implementations
of poker primitives (cards, deck, hand evaulations, etc)....





Runtime dependencies:   

     1.  postgres
     2.  Centrifugo
     3.  keycloak





Setup overview:

0.   Basic config env vars:

    LOCAL_IP
    POKERLIB_HANDEVAL_LOG

1.  Install postgres

     env vars:
       
      POSTGRES_USER
      POSTGRES_PASSWORD
      POKER_DB_HOST
      
2.  Install Centrifugo

     env vars:
       
      CENTRIFUGO_SERVICE
      CENTRIFUGO_HMAC_SECRET
      CENTRIFUGO_CLIENT_URL
      CENTRIFUGO_KEY
      
 
3.  Install Keycloak

     env vars:
     
     KEYCLOAK_CLIENT_ID
     KEYCLOAK_CLIENT_SECRET
     AUTH_PROVIDER_DOMAIN
     POKER_CALLBACK_URL

4.  Configure Keycloak with an IDP (google)


