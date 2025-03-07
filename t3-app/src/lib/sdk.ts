// import { env } from "@/env";
// import { Client, Databases, ExecutionMethod, Functions, Query } from "appwrite";
// import type { CheckList, HostList, Network, StatusList } from "./types";

// const client = new Client()
//   .setEndpoint(env.NEXT_PUBLIC_APPWRITE_ENDPOINT) // Your API Endpoint
//   .setProject(env.NEXT_PUBLIC_APPWRITE_PROJECTID); // Your project ID

// const databases = new Databases(client);
// const serverless = new Functions(client);

// class Appwrite {
//   async getOnlineHosts(network: Network, netAddress: string) {
//     try {
//       const hosts = (await databases.listDocuments(
//         network === "mainnet"
//           ? env.NEXT_PUBLIC_APPWRITE_DBMAIN
//           : env.NEXT_PUBLIC_APPWRITE_DBZEN,
//         network === "mainnet"
//           ? env.NEXT_PUBLIC_APPWRITE_COLHOSTS_MAIN
//           : env.NEXT_PUBLIC_APPWRITE_COLHOSTS_ZEN,
//         [
//           Query.contains("netAddress", netAddress),
//           Query.limit(200),
//           Query.equal("online", true),
//         ],
//       )) as HostList;

//       if (hosts.total !== hosts.documents.length) {
//         // round upwards to next whole integer
//         const needToGrab = Math.ceil(
//           (hosts.total - hosts.documents.length) / 200,
//         );

//         for (let i = 0; i < needToGrab; i++) {
//           const newHosts = (await databases.listDocuments(
//             network === "mainnet"
//               ? env.NEXT_PUBLIC_APPWRITE_DBMAIN
//               : env.NEXT_PUBLIC_APPWRITE_DBZEN,
//             network === "mainnet"
//               ? env.NEXT_PUBLIC_APPWRITE_COLHOSTS_MAIN
//               : env.NEXT_PUBLIC_APPWRITE_COLHOSTS_ZEN,
//             [
//               Query.contains("netAddress", netAddress),
//               Query.limit(200),
//               Query.equal("online", true),
//               Query.offset(hosts.documents.length),
//             ],
//           )) as HostList;
//           hosts.documents.push(...newHosts.documents);
//         }
//       }
//       return hosts;
//     } catch {
//       return null;
//     }
//   }

//   async getAllHosts(network: Network, netAddress: string) {
//     try {
//       if (netAddress == "") {
//         return null;
//       }
//       const hosts = (await databases.listDocuments(
//         network === "mainnet"
//           ? env.NEXT_PUBLIC_APPWRITE_DBMAIN
//           : env.NEXT_PUBLIC_APPWRITE_DBZEN,
//         network === "mainnet"
//           ? env.NEXT_PUBLIC_APPWRITE_COLHOSTS_MAIN
//           : env.NEXT_PUBLIC_APPWRITE_COLHOSTS_ZEN,
//         [Query.contains("netAddress", netAddress), Query.limit(200)],
//       )) as HostList;

//       if (hosts.total !== hosts.documents.length) {
//         // round upwards to next whole integer
//         const needToGrab = Math.ceil(
//           (hosts.total - hosts.documents.length) / 200,
//         );

//         for (let i = 0; i < needToGrab; i++) {
//           const newHosts = (await databases.listDocuments(
//             network === "mainnet"
//               ? env.NEXT_PUBLIC_APPWRITE_DBMAIN
//               : env.NEXT_PUBLIC_APPWRITE_DBZEN,
//             network === "mainnet"
//               ? env.NEXT_PUBLIC_APPWRITE_COLHOSTS_MAIN
//               : env.NEXT_PUBLIC_APPWRITE_COLHOSTS_ZEN,
//             [
//               Query.contains("netAddress", netAddress),
//               Query.limit(200),
//               Query.offset(hosts.documents.length),
//             ],
//           )) as HostList;
//           hosts.documents.push(...newHosts.documents);
//         }
//       }
//       return hosts;
//     } catch {
//       return null;
//     }
//   }

//   async getHosts(network: Network, netAddress: string, online: boolean) {
//     try {
//       if (online) {
//         return await this.getOnlineHosts(network, netAddress);
//       }
//       return await this.getAllHosts(network, netAddress);
//     } catch {
//       return null;
//     }
//   }

//   async getStatus(network: Network) {
//     try {
//       const status = (await databases.listDocuments(
//         network === "mainnet"
//           ? env.NEXT_PUBLIC_APPWRITE_DBMAIN
//           : env.NEXT_PUBLIC_APPWRITE_DBZEN,
//         network === "mainnet"
//           ? env.NEXT_PUBLIC_APPWRITE_COLSTATUS_MAIN
//           : env.NEXT_PUBLIC_APPWRITE_COLSTATUS_ZEN,
//       )) as StatusList;
//       return status.documents[0];
//     } catch {
//       return null;
//     }
//   }

//   async sendHostEmail(network: Network, publicKey: string, email: string) {
//     try {
//       const response = await serverless.createExecution(
//         "67b491080036e094a8cc",
//         JSON.stringify({ network, publicKey, email }),
//         false,
//         undefined,
//         ExecutionMethod.POST,
//       );
//       if (response.responseStatusCode === 201) {
//         return true;
//       }
//       // console.log(response.responseBody);
//       return false;
//     } catch {
//       return false;
//     }
//   }

//   async verifyOtp(
//     network: Network,
//     secret: string,
//     publicKey: string,
//     email: string,
//     expire: string,
//   ) {
//     try {
//       const response = await serverless.createExecution(
//         "67b491080036e094a8cc",
//         JSON.stringify({ network, secret, publicKey, email, expire }),
//         false,
//         undefined,
//         ExecutionMethod.PUT,
//       );
//       if (response.responseStatusCode === 200) {
//         return JSON.parse(response.responseBody) as {
//           message: string;
//           publicKey: string;
//           email: string;
//           address: string;
//         };
//       } else {
//         throw new Error(response.responseBody);
//       }
//     } catch {
//       return null;
//     }
//   }
//   async checkServerless() {
//     try {
//       const response = await serverless.createExecution(
//         "67b491080036e094a8cc",
//         JSON.stringify({}),
//         false,
//         undefined,
//         ExecutionMethod.GET,
//       );
//       return response.responseStatusCode === 200;
//     } catch {
//       return false;
//     }
//   }

//   async scanHost(network: Network, hostId: string) {
//     try {
//       return (await databases.listDocuments(
//         network === "mainnet"
//           ? env.NEXT_PUBLIC_APPWRITE_DBMAIN
//           : env.NEXT_PUBLIC_APPWRITE_DBZEN,
//         network === "mainnet"
//           ? env.NEXT_PUBLIC_APPWRITE_COLCHECK_MAIN
//           : env.NEXT_PUBLIC_APPWRITE_COLCHECK_ZEN,
//         [Query.equal("hostId", hostId)],
//       )) as CheckList;
//     } catch (error) {
//       console.log(error);
//       return null;
//     }
//   }
// }
// const sdk = new Appwrite();
// export default sdk;
