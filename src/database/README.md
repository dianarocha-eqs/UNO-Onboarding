# GO API

## Database Access

To connect to the Azure SQL Database for this project, use the following credentials:

### **Connection Details**
- **Server**: `uno-onboarding.database.windows.net`
- **Database**: `uno-onboarding`
- **User ID**: `eqs-digital`
- **Password**: `B3N5GA!QpgCX^&@Bt+pXhTY6NcZD7NWw`
- **Port**: `1433`

### **Connection String**
For applications or tools that require a connection string, use the following:

```plaintext
Server=tcp:uno-onboarding.database.windows.net,1433;
Initial Catalog=uno-onboarding;
Persist Security Info=False;
User ID=eqs-digital;
Password=B3N5GA!QpgCX^&@Bt+pXhTY6NcZD7NWw;
MultipleActiveResultSets=False;
Encrypt=True;
TrustServerCertificate=False;
Connection Timeout=30;
```

### **Security Note**
- The firewall currently allows access from all IP addresses. Ensure you connect securely and avoid sharing credentials publicly.


