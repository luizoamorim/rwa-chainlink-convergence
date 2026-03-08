# Tenderly Setup

This project uses **Tenderly Virtual Testnets** as the blockchain environment.

Official website:

https://tenderly.co/

---

# 1. Create an Account

Go to:

https://tenderly.co/

Click **Go to Dashboard**.

Login or register.

You can use:

* Google
* GitHub

---

# 2. Create an Organization

Navigate to:

```
Members → Create Organization
```

Use the **Free Trial plan**.

Example name:

```
hackathonteam
```

---

# 3. Create a Project

Inside your organization:

Create a new project.

Example name:

```
hackthonproject
```

---

# 4. Create a Virtual Testnet

Inside the project dashboard:

Go to:

```
Virtual Testnets
```

Click:

```
Create Virtual Testnet
```

---

# 5. Copy RPC URLs

After creating the testnet you will see:

```
Admin RPC (HTTPS)
Admin RPC (WebSocket)
Public RPC (Explorer)
```

Copy both values.

---

# 6. Configure Environment Variables

Add the RPC URL to your `.env` file.

⚠️ **IMPORTANT**

For this step you must use the **Admin RPC (HTTPS)**.

The **Admin RPC is private** and allows contract deployment and blockchain writes.

Example:

```
TENDERLY_RPC_URL=https://rpc.tenderly.co/...
```

---

# 7. Configure CRE Project

⚠️ **IMPORTANT**

The **CRE workflow also requires the Admin RPC (HTTPS)**.

Open:

```
project.yaml
```

Set the staging RPC URL using the same value.

Example:

```
rpc_url: https://rpc.tenderly.co/...
```

---

## 8. Open the Tenderly Explorer

Add the PUBLIC RPC URL to your `.env.local` file in the `frontend` project.

⚠️ **IMPORTANT**

For this step you must use the **Public RPC URL**, not the Admin RPC.

The Public RPC allows you to access the **Tenderly transaction explorer**.

To get it:

1. Copy the **Public RPC URL**
2. Paste it into a new browser tab
3. Open the URL
4. Copy the **full URL from the browser**

It will look like this: https://dashboard.tenderly.co/...

---

# Setup Completed

Tenderly is now ready.

You can now run the project deployment.
