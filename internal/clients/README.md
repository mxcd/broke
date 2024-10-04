## Keycloak API strategy

1. Request groups from kc `/admin/realms/{realm}/groups/{id}`
2. Request roles of groups from kc `/admin/realms/{realm}/groups/{id}`

### Case 1: Only kc groups are requested

3. Get members of requested groups from kc `/admin/realms/{realm}/groups/{id}/members`
4. Map kc groups and roles to members

### Case 2: Some kc roles are requested

3. Get groups for every requested roles
4. Get all requested groups and merge with groups from step 3
4. Get all groups of users from kc `/admin/realms/{realm}/users/{id}/groups`
5. Get all roles of users from kc `/admin/realms/{realm}/users/{id}/role-mappings/realm`
6. Map kc groups and roles to users
