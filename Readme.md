 
 * Note:   This is currently not in a runnable state, it is midway through
           a refactor replacing centrifugo with NATS and keycloack with auth0...



 The poker app assumes the availability these third party tools.
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


