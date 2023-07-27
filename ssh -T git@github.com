[33mcommit bbdb576c0cc5bf5f831d8169bd6f744b2b723597[m[33m ([m[1;36mHEAD -> [m[1;32mPSL-898_WNSetupMesh[m[33m, [m[1;31morigin/PSL-898_WNSetupMesh[m[33m)[m
Author: j-rafique <mejrafique@gmail.com>
Date:   Thu Jul 27 13:27:30 2023 +0500

    [PSL-898] implement revised setup mesh logic on WN

[33mcommit f283b93eec1f81ee78ffce3bd685a20e28a2f1ff[m[33m ([m[1;31morigin/master[m[33m, [m[1;31morigin/HEAD[m[33m, [m[1;32mmaster[m[33m)[m
Merge: b6511b99 4ed2300a
Author: J Bilal rafique <113895287+j-rafique@users.noreply.github.com>
Date:   Wed Jul 26 12:14:38 2023 -0700

    Merge pull request #623 from pastelnetwork/PSL-897_getTopMNsResponseFromSN
    
    [PSL-897] implement the proto and SN logic to get MN-top list

[33mcommit 4ed2300a9bfe47f90f1ca0034a58b4bb680d6fd8[m[33m ([m[1;31morigin/PSL-897_getTopMNsResponseFromSN[m[33m, [m[1;32mPSL-897_getTopMNsResponseFromSN[m[33m)[m
Author: j-rafique <mejrafique@gmail.com>
Date:   Wed Jul 26 09:30:12 2023 -0700

    [PSL-897] implement the proto and SN logic to get MN-top list

[33mcommit b6511b992aa4a8ab07b2b57e173dc541bf815e87[m
Merge: 0b1bfc13 0d22fe8c
Author: J Bilal rafique <113895287+j-rafique@users.noreply.github.com>
Date:   Wed Jul 26 19:31:23 2023 +0500

    Merge pull request #622 from pastelnetwork/PSL-896_defineRequestProtoByWN
    
    [PSL-896] implement grpc request contract by defining proto on WN side

[33mcommit 0d22fe8c00e515e063532e14c2d242333275908f[m[33m ([m[1;31morigin/PSL-896_defineRequestProtoByWN[m[33m, [m[1;32mPSL-896_defineRequestProtoByWN[m[33m)[m
Author: j-rafique <mejrafique@gmail.com>
Date:   Wed Jul 26 07:08:13 2023 -0700

    [PSL-896] implement grpc request contract by defining proto on WN side

[33mcommit 0b1bfc1313afb54b0fcb19cf684d191e708b6222[m[33m ([m[1;33mtag: v1.2.17-beta[m[33m)[m
Merge: bbf6423b db03e41d
Author: Matee ullah Malik <46045452+mateeullahmalik@users.noreply.github.com>
Date:   Tue Jul 25 05:13:17 2023 +0500

    Merge pull request #621 from pastelnetwork/PSL-891_connClose
    
    fix conn close

[33mcommit db03e41d12ef98ab43c62908141e2922b7eb7657[m
Author: Matee <matee.u@hotmail.com>
Date:   Tue Jul 25 04:58:09 2023 +0500

    fix conn close

[33mcommit bbf6423b5847863fdbff9244e1b1c1d2ae2306ee[m[33m ([m[1;33mtag: v1.2.16-beta[m[33m)[m
Merge: 7eb5bec1 0445815f
Author: Matee ullah Malik <46045452+mateeullahmalik@users.noreply.github.com>
Date:   Mon Jul 24 21:58:43 2023 +0500

    Merge pull request #620 from pastelnetwork/PSL-891_connClose
    
    Psl 891 conn close

[33mcommit 0445815fe8041c9df9cc9ba304158a26a1431264[m
Author: Matee <matee.u@hotmail.com>
Date:   Mon Jul 24 21:34:44 2023 +0500

    searching node

[33mcommit 42b2a4359a1d0a335a58ba68925d8c5360ea4488[m
Author: Matee <matee.u@hotmail.com>
Date:   Mon Jul 24 19:16:13 2023 +0500

    p2p replication fix

[33mcommit a36efbaca9c1d28c6a4d3f4634077858cb72e7c5[m
Author: Matee <matee.u@hotmail.com>
Date:   Sun Jul 23 17:30:43 2023 +0500

    tmp

[33mcommit 7eb5bec1b469e6a493858db7959bd3b2fb09f71c[m
Merge: 24ae5c4e fce950a0
Author: Matee ullah Malik <46045452+mateeullahmalik@users.noreply.github.com>
Date:   Sat Jul 22 05:44:01 2023 +0500

    Merge pull request #619 from pastelnetwork/fixLogAndBulkCascadeRegScript
    
    fix mesh handler close connection log & cascade reg script

[33mcommit fce950a06fd19018b3ec692ce497a22241fba6ad[m[33m ([m[1;31morigin/fixLogAndBulkCascadeRegScript[m[33m, [m[1;32mfixLogAndBulkCascadeRegScript[m[33m)[m
Author: j-rafique <mejrafique@gmail.com>
Date:   Fri Jul 21 07:45:27 2023 -0700

    fix mesh handler close connection log & cascade reg script

[33mcommit 24ae5c4e692212fc28ef774a6ef1544cdcaf072f[m
Merge: bc332622 e4953275
Author: J Bilal rafique <113895287+j-rafique@users.noreply.github.com>
Date:   Thu Jul 20 23:16:34 2023 -0700

    Merge pull request #618 from pastelnetwork/PSL-890_fixGetDDServerStatsLogs
    
    [PSL-890] fix get dd-server stats logs & bulk sense reg script

[33mcommit e49532757575e4284762b6b57265f51dadce1f6a[m[33m ([m[1;31morigin/PSL-890_fixGetDDServerStatsLogs[m[33m, [m[1;32mPSL-890_fixGetDDServerStatsLogs[m[33m)[m
Author: j-rafique <mejrafique@gmail.com>
Date:   Thu Jul 20 23:29:15 2023 +0500

    [PSL-890] fix get dd-server stats logs & bulk sense reg script

[33mcommit bc332622b075cb2a26e74512216cc9e194efc13b[m
Merge: d32fc9ea 34c20c06
Author: Matee ullah Malik <46045452+mateeullahmalik@users.noreply.github.com>
Date:   Wed Jul 19 16:28:20 2023 +0500

    Merge pull request #617 from pastelnetwork/PSL-876_modifyIterateFixReplication
    
    Psl 876 modify iterate fix replication

[33mcommit d32fc9ea3aca65f10dba55b6f6116a01dbf42880[m
Merge: 07ab70c4 f31b8d92
Author: J Bilal rafique <113895287+j-rafique@users.noreply.github.com>
Date:   Wed Jul 19 04:26:04 2023 -0700

    Merge pull request #615 from pastelnetwork/PSL-881_DDServerQueueChanges
    
    [PSL-881] implement changes to check dd-server availability for the SN

[33mcommit 34c20c069d6c4f5bdd0efaf23b2a6f45dc51c5b2[m
Author: Matee <matee.u@hotmail.com>
Date:   Wed Jul 19 15:36:51 2023 +0500

    fix replication adjust keys on node leaving case

[33mcommit f2979c95d2c58655824c176939b2a320fc334187[m
Author: Matee <matee.u@hotmail.com>
Date:   Wed Jul 19 15:34:20 2023 +0500

    fix replication adjust keys on node leaving case

[33mcommit 829516460026578df1bc6a296a83ff852d3263eb[m
Author: Matee <matee.u@hotmail.com>
Date:   Wed Jul 19 14:22:40 2023 +0500

    fix

[33mcommit f31b8d92ead75c1e5f4e9a8c316140379bfb718b[m[33m ([m[1;31morigin/PSL-881_DDServerQueueChanges[m[33m, [m[1;32mPSL-881_DDServerQueueChanges[m[33m)[m
Author: j-rafique <mejrafique@gmail.com>
Date:   Wed Jul 19 02:09:42 2023 -0700

    [PSL-881] implement changes to check dd-server availability for the SN

[33mcommit 07ab70c4fdd4ad55c2977162a25e776e85e98d0c[m
Merge: ee78ffeb eae2e9c3
Author: a-ok123 <54385956+a-ok123@users.noreply.github.com>
Date:   Tue Jul 18 02:32:08 2023 -0400

    Merge pull request #616 from pastelnetwork/cli-fix
    
    [PSL-765] Fix for help not printing flags

[33mcommit eae2e9c30f45f001b9b0c5073bd280981002b81e[m
Author: a-ok123 <alexey@pastel.network>
Date:   Tue Jul 18 02:09:59 2023 -0400

    Fix for help not printing flags

[33mcommit ee78ffebe453161f4e72a751b9f7b4ce2a380df6[m[33m ([m[1;33mtag: v1.2.15-beta[m[33m)[m
Merge: 80aedd65 d0c9e000
Author: Matee ullah Malik <46045452+mateeullahmalik@users.noreply.github.com>
Date:   Tue Jul 18 02:48:45 2023 +0500

    Merge pull request #614 from pastelnetwork/PSL-876_modifyIterateFixReplication
    
    [PSL-848] save p2p file type with data, [PSL-876] modify iterate, fix…

[33mcommit d0c9e00038def7802498f03584d7fbd433e2b66c[m
Author: Matee <matee.u@hotmail.com>
Date:   Tue Jul 18 02:34:16 2023 +0500

    fix

[33mcommit 53ec48c7ae636f966450b47cacca45ecb8f76819[m
Author: Matee <matee.u@hotmail.com>
Date:   Mon Jul 17 17:38:28 2023 +0500

    fix replication

[33mcommit e4d0cddd278681118d62ca9ed421a9d4b4b290b8[m
Author: Matee <matee.u@hotmail.com>
Date:   Mon Jul 17 16:20:40 2023 +0500

    tmp

[33mcommit 63c379b8ac3f754a1a2d1aa30b8f16d5341893d1[m
Author: Matee <matee.u@hotmail.com>
Date:   Mon Jul 17 16:15:35 2023 +0500

    tmp

[33mcommit 6843a72ce9a789fcf9fe6526c87aa9162215ee95[m
Author: Matee <matee.u@hotmail.com>
Date:   Mon Jul 17 14:03:48 2023 +0500

    replication worker

[33mcommit 7a952f7b969fc5c96932830c7b65a8a442a8e35f[m
Author: Matee <matee.u@hotmail.com>
Date:   Sun Jul 16 06:18:08 2023 +0500

    fix lint

[33mcommit 70ee8489d6e29694380760829cc6e0ef20ce994f[m
Author: Matee <matee.u@hotmail.com>
Date:   Sun Jul 16 03:14:56 2023 +0500

    close conn
    
    Replication mechanism redesign

[33mcommit d800b2d97f143ee2a7dfc9ac90ea32d296ee2c36[m
Author: Matee <matee.u@hotmail.com>
Date:   Sat Jul 15 04:52:28 2023 +0500

    improve conn pool

[33mcommit a3528975bcba80895a1cbc86ea708f3e00aaa265[m[33m ([m[1;31morigin/PSL-876_modifyIterateFixReplication[m[33m)[m
Author: Matee <matee.u@hotmail.com>
Date:   Tue Jul 11 02:06:02 2023 +0500

    fix lint

[33mcommit 