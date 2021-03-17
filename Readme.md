# ec2price

This is a simple cli tool for quickly comparing ec2 instance types.

## Sample output

```
# list instance types
$ ./ec2price
           type        mem   vcpu            disk   dsk mfg    hourly    annual annual-reserved
       t4g.nano        0.5      2        EBS only     0 arm    0.0042     36.79 26.28
       t3a.nano        0.5      2        EBS only     0 amd    0.0047     41.17 29.78
        t3.nano        0.5      2        EBS only     0 int    0.0052     45.55 32.41
        t2.nano        0.5      1        EBS only     0 int    0.0058     50.81 35.92
      t4g.micro        1.0      2        EBS only     0 arm    0.0084     73.58 53.44
      t3a.micro        1.0      2        EBS only     0 amd    0.0094     82.34 59.57
       t3.micro        1.0      2        EBS only     0 int    0.0104     91.10 65.70
       t2.micro        1.0      1        EBS only     0 int    0.0116    101.62 72.71
      t4g.small        2.0      2        EBS only     0 arm    0.0168    147.17 106.00
      t3a.small        2.0      2        EBS only     0 amd    0.0188    164.69 119.14
       t1.micro        0.6      1        EBS only     0 unk    0.0200    175.20 0.00
       t3.small        2.0      2        EBS only     0 int    0.0208    182.21 131.40
       t2.small        2.0      1        EBS only     0 int    0.0230    201.48 144.54
      a1.medium        2.0      1        EBS only     0 arm    0.0255    223.38 162.06
     t4g.medium        4.0      2        EBS only     0 arm    0.0336    294.34 211.99
     c6g.medium        2.0      1        EBS only     0 arm    0.0340    297.84 217.25
     t3a.medium        4.0      2        EBS only     0 amd    0.0376    329.38 237.40
    c6gd.medium        2.0      1     1 x 59 NVMe    59 arm    0.0384    336.38 246.16
     m6g.medium        4.0      1        EBS only     0 arm    0.0385    337.26 247.91
      t3.medium        4.0      2        EBS only     0 int    0.0416    364.42 262.80
    c6gn.medium        2.0      1        EBS only     0 arm    0.0432    378.43 286.63
    m6gd.medium        4.0      1     1 x 59 NVMe    59 arm    0.0452    395.95 290.83
      t2.medium        4.0      2        EBS only     0 int    0.0464    406.46 289.96
     r6g.medium        8.0      1        EBS only     0 arm    0.0504    441.50 319.74
       a1.large        4.0      2        EBS only     0 arm    0.0510    446.76 323.24
    r6gd.medium        8.0      1     1 x 59 NVMe    59 arm    0.0576    504.58 367.92
      m3.medium        3.8      1           1 x 4     4 int    0.0670    586.92 438.00
      t4g.large        8.0      2        EBS only     0 arm    0.0672    588.67 424.86
      c6g.large        4.0      2        EBS only     0 arm    0.0680    595.68 434.50
      t3a.large        8.0      2        EBS only     0 amd    0.0752    658.75 474.79
     c6gd.large        4.0      2    1 x 118 NVMe   118 arm    0.0768    672.77 493.19
      m6g.large        8.0      2        EBS only     0 arm    0.0770    674.52 496.69
      c5a.large        4.0      2        EBS only     0 amd    0.0770    674.52 490.56
       t3.large        8.0      2        EBS only     0 int    0.0832    728.83 525.60
       c5.large        4.0      2        EBS only     0 int    0.0850    744.60 543.12
     c5ad.large        4.0      2     1 x 75 NVMe    75 amd    0.0860    753.36 543.12
      m5a.large        8.0      2        EBS only     0 amd    0.0860    753.36 551.88
     c6gn.large        4.0      2        EBS only     0 arm    0.0864    756.86 573.25
      m1.medium        3.8      1         1 x 410   410 int    0.0870    762.12 0.00
     m6gd.large        8.0      2    1 x 118 NVMe   118 arm    0.0904    791.90 581.66
       t2.large        8.0      2        EBS only     0 int    0.0928    812.93 579.04
       m5.large        8.0      2        EBS only     0 int    0.0960    840.96 621.96
      c5d.large        4.0      2     1 x 50 NVMe    50 int    0.0960    840.96 613.20
       m4.large        8.0      2        EBS only     0 int    0.1000    876.00 623.71
       c4.large        3.8      2        EBS only     0 int    0.1000    876.00 630.72
      r6g.large       16.0      2        EBS only     0 arm    0.1008    883.01 640.36
      a1.xlarge        8.0      4        EBS only     0 arm    0.1020    893.52 647.36
     m5ad.large        8.0      2     1 x 75 NVMe    75 amd    0.1030    902.28 665.76
       c3.large        3.8      2          2 x 16    32 int    0.1050    919.80 665.76
      c5n.large        5.2      2        EBS only     0 int    0.1080    946.08 683.28
      r5a.large       16.0      2        EBS only     0 amd    0.1130    989.88 718.32
      m5d.large        8.0      2     1 x 75 NVMe    75 int    0.1130    989.88 727.08
     r6gd.large       16.0      2    1 x 118 NVMe   118 arm    0.1152   1009.15 735.84
      m5n.large        8.0      2        EBS only     0 int    0.1190   1042.44 762.12
       r5.large       16.0      2        EBS only     0 int    0.1260   1103.76 797.16
      c1.medium        1.7      2         1 x 350   350 int    0.1300   1138.80 0.00
     r5ad.large       16.0      2     1 x 75 NVMe    75 amd    0.1310   1147.56 832.20
       m3.large        7.5      2          1 x 32    32 int    0.1330   1165.08 867.24
       r4.large       15.2      2        EBS only     0 int    0.1330   1165.08 846.22
     t4g.xlarge       16.0      4        EBS only     0 arm    0.1344   1177.34 848.84
     m5dn.large        8.0      2     1 x 75 NVMe    75 int    0.1360   1191.36 867.24
     c6g.xlarge        8.0      4        EBS only     0 arm    0.1360   1191.36 868.12
      r5d.large       16.0      2     1 x 75 NVMe    75 int    0.1440   1261.44 919.80
      r5b.large       16.0      2        EBS only     0 int    0.1490   1305.24 945.65
      r5n.large       16.0      2        EBS only     0 int    0.1490   1305.24 946.08
     t3a.xlarge       16.0      4        EBS only     0 amd    0.1504   1317.50 950.46
    c6gd.xlarge        8.0      4    1 x 237 NVMe   237 arm    0.1536   1345.54 986.38
     m6g.xlarge       16.0      4        EBS only     0 arm    0.1540   1349.04 993.38
     c5a.xlarge        8.0      4        EBS only     0 amd    0.1540   1349.04 981.12
       i3.large       15.2      2    1 x 475 NVMe   475 int    0.1560   1366.56 1077.48
     m5zn.large        8.0      2        EBS only     0 int    0.1652   1447.15 1094.12
       r3.large       15.2      2          1 x 32    32 int    0.1660   1454.16 946.08
      t3.xlarge       16.0      4        EBS only     0 int    0.1664   1457.66 1051.20
     r5dn.large       16.0      2     1 x 75 NVMe    75 int    0.1670   1462.92 1059.96
      c5.xlarge        8.0      4        EBS only     0 int    0.1700   1489.20 1077.48
     m5a.xlarge       16.0      4        EBS only     0 amd    0.1720   1506.72 1103.76
    c5ad.xlarge        8.0      4    1 x 150 NVMe   150 amd    0.1720   1506.72 1086.24
    c6gn.xlarge        8.0      4        EBS only     0 arm    0.1728   1513.73 1146.42
       m1.large        7.5      2         2 x 420   840 int    0.1750   1533.00 0.00
    m6gd.xlarge       16.0      4    1 x 237 NVMe   237 arm    0.1808   1583.81 1162.45
      t2.xlarge       16.0      4        EBS only     0 int    0.1856   1625.86 1158.07
      z1d.large       16.0      2     1 x 75 NVMe    75 int    0.1860   1629.36 1182.60
      m5.xlarge       16.0      4        EBS only     0 int    0.1920   1681.92 1235.16
     c5d.xlarge        8.0      4    1 x 100 NVMe   100 int    0.1920   1681.92 1217.64
      c4.xlarge        7.5      4        EBS only     0 int    0.1990   1743.24 1270.20
      m4.xlarge       16.0      4        EBS only     0 int    0.2000   1752.00 1248.30
     r6g.xlarge       32.0      4        EBS only     0 arm    0.2016   1766.02 1280.71
     a1.2xlarge       16.0      8        EBS only     0 arm    0.2040   1787.04 1294.73
    m5ad.xlarge       16.0      4    1 x 150 NVMe   150 amd    0.2060   1804.56 1322.76
      c3.xlarge        7.5      4          2 x 40    80 int    0.2100   1839.60 1340.28
     c5n.xlarge       10.5      4        EBS only     0 int    0.2160   1892.16 1375.32
     m5d.xlarge       16.0      4    1 x 150 NVMe   150 int    0.2260   1979.76 1445.40
     r5a.xlarge       32.0      4        EBS only     0 amd    0.2260   1979.76 1436.64
     i3en.large       16.0      2   1 x 1250 NVMe  1250 int    0.2260   1979.76 1559.28
    r6gd.xlarge       32.0      4    1 x 237 NVMe   237 arm    0.2304   2018.30 1471.68
     m5n.xlarge       16.0      4        EBS only     0 int    0.2380   2084.88 1524.24
      m2.xlarge       17.1      2         1 x 420   420 int    0.2450   2146.20 0.00
      r5.xlarge       32.0      4        EBS only     0 int    0.2520   2207.52 1603.08
    r5ad.xlarge       32.0      4    1 x 150 NVMe   150 amd    0.2620   2295.12 1673.16
      r4.xlarge       30.5      4        EBS only     0 int    0.2660   2330.16 1692.43
      m3.xlarge       15.0      4          2 x 40    80 int    0.2660   2330.16 1743.24
    t4g.2xlarge       32.0      8        EBS only     0 arm    0.2688   2354.69 1697.69
    c6g.2xlarge       16.0      8        EBS only     0 arm    0.2720   2382.72 1736.23
    m5dn.xlarge       16.0      4    1 x 150 NVMe   150 int    0.2720   2382.72 1743.24
     r5d.xlarge       32.0      4    1 x 150 NVMe   150 int    0.2880   2522.88 1830.84
     r5n.xlarge       32.0      4        EBS only     0 int    0.2980   2610.48 1892.16
     r5b.xlarge       32.0      4        EBS only     0 int    0.2980   2610.48 1891.29
    t3a.2xlarge       32.0      8        EBS only     0 amd    0.3008   2635.01 1900.92
   c6gd.2xlarge       16.0      8    1 x 475 NVMe   475 arm    0.3072   2691.07 1971.88
    c5a.2xlarge       16.0      8        EBS only     0 amd    0.3080   2698.08 1962.24
    m6g.2xlarge       32.0      8        EBS only     0 arm    0.3080   2698.08 1986.77
      i3.xlarge       30.5      4    1 x 950 NVMe   950 int    0.3120   2733.12 2154.96
    m5zn.xlarge       16.0      4        EBS only     0 int    0.3303   2893.43 2187.37
     t3.2xlarge       32.0      8        EBS only     0 int    0.3328   2915.33 2101.52
      r3.xlarge       30.5      4          1 x 80    80 int    0.3330   2917.08 1900.92
    r5dn.xlarge       32.0      4    1 x 150 NVMe   150 int    0.3340   2925.84 2119.92
     c5.2xlarge       16.0      8        EBS only     0 int    0.3400   2978.40 2154.96
   c5ad.2xlarge       16.0      8    1 x 300 NVMe   300 amd    0.3440   3013.44 2181.24
    m5a.2xlarge       32.0      8        EBS only     0 amd    0.3440   3013.44 2207.52
   c6gn.2xlarge       16.0      8        EBS only     0 arm    0.3456   3027.46 2292.93
      m1.xlarge       15.0      4         4 x 420  1680 int    0.3500   3066.00 0.00
   m6gd.2xlarge       32.0      8    1 x 475 NVMe   475 arm    0.3616   3167.62 2324.90
    inf1.xlarge        8.0      4        EBS only     0 int    0.3680   3223.68 2338.92
     t2.2xlarge       32.0      8        EBS only     0 int    0.3712   3251.71 2317.02
     z1d.xlarge       32.0      4    1 x 150 NVMe   150 int    0.3720   3258.72 2365.20
    c5d.2xlarge       16.0      8    1 x 200 NVMe   200 int    0.3840   3363.84 2444.04
     m5.2xlarge       32.0      8        EBS only     0 int    0.3840   3363.84 2470.32
     c4.2xlarge       15.0      8        EBS only     0 int    0.3980   3486.48 2540.40
     m4.2xlarge       32.0      8        EBS only     0 int    0.4000   3504.00 2496.60
    r6g.2xlarge       64.0      8        EBS only     0 arm    0.4032   3532.03 2560.55
     a1.4xlarge       32.0     16        EBS only     0 arm    0.4080   3574.08 2589.46
       a1.metal       32.0     16        EBS only     0 arm    0.4080   3574.08 2592.96
   m5ad.2xlarge       32.0      8    1 x 300 NVMe   300 amd    0.4120   3609.12 2654.28
     c3.2xlarge       15.0      8          2 x 80   160 int    0.4200   3679.20 2671.80
    c5n.2xlarge       21.0      8        EBS only     0 int    0.4320   3784.32 2741.88
    i3en.xlarge       32.0      4   1 x 2500 NVMe  2500 int    0.4520   3959.52 3118.56
    r5a.2xlarge       64.0      8        EBS only     0 amd    0.4520   3959.52 2873.28
    m5d.2xlarge       32.0      8    1 x 300 NVMe   300 int    0.4520   3959.52 2899.56
   r6gd.2xlarge       64.0      8    1 x 475 NVMe   475 arm    0.4608   4036.61 2943.36
     h1.2xlarge       32.0      8    1 x 2000 HDD  2000 int    0.4680   4099.68 3206.16
    m5n.2xlarge       32.0      8        EBS only     0 int    0.4760   4169.76 3039.72
     m2.2xlarge       34.2      4         1 x 850   850 int    0.4900   4292.40 0.00
      d3.xlarge       32.0      4    3 x 2000 HDD  6000 int    0.4990   4371.24 3311.28
     r5.2xlarge       64.0      8        EBS only     0 int    0.5040   4415.04 3197.40
      c1.xlarge        7.0      8         4 x 420  1680 int    0.5200   4555.20 0.00
   r5ad.2xlarge       64.0      8    1 x 300 NVMe   300 amd    0.5240   4590.24 3337.56
    g4dn.xlarge       16.0      4     125 GB NVMe     0 int    0.5260   4607.76 3337.56
    d3en.xlarge       16.0      4   2 x 14000 HDD 28000 int    0.5260   4607.76 3477.72
     r4.2xlarge       61.0      8        EBS only     0 int    0.5320   4660.32 3384.86
     m3.2xlarge       30.0      8          2 x 80   160 int    0.5320   4660.32 3477.72
   m5dn.2xlarge       32.0      8    1 x 300 NVMe   300 int    0.5440   4765.44 3486.48
    c6g.4xlarge       32.0     16        EBS only     0 arm    0.5440   4765.44 3473.34
    r5d.2xlarge       64.0      8    1 x 300 NVMe   300 int    0.5760   5045.76 3661.68
   inf1.2xlarge       16.0      8        EBS only     0 int    0.5840   5115.84 3705.48
    r5n.2xlarge       64.0      8        EBS only     0 int    0.5960   5220.96 3784.32
    r5b.2xlarge       64.0      8        EBS only     0 int    0.5960   5220.96 3782.59
   c6gd.4xlarge       32.0     16    1 x 950 NVMe   950 arm    0.6144   5382.14 3943.75
    m6g.4xlarge       64.0     16        EBS only     0 arm    0.6160   5396.16 3972.66
    c5a.4xlarge       32.0     16        EBS only     0 amd    0.6160   5396.16 3924.48
     i3.2xlarge       61.0      8   1 x 1900 NVMe  1900 int    0.6240   5466.24 4318.68
     g2.2xlarge       15.0      8          1 x 60    60 int    0.6500   5694.00 4283.64
   m5zn.2xlarge       32.0      8        EBS only     0 int    0.6607   5787.73 4375.62
     r3.2xlarge       61.0      8         1 x 160   160 int    0.6650   5825.40 3793.08
   r5dn.2xlarge       64.0      8    1 x 300 NVMe   300 int    0.6680   5851.68 4239.84
     c5.4xlarge       32.0     16        EBS only     0 int    0.6800   5956.80 4318.68
   c5ad.4xlarge       32.0     16    2 x 300 NVMe   600 amd    0.6880   6026.88 4362.48
    m5a.4xlarge       64.0     16        EBS only     0 amd    0.6880   6026.88 4415.04
      d2.xlarge       30.5      4    3 x 2000 HDD  6000 int    0.6900   6044.40 3652.92
   c6gn.4xlarge       32.0     16        EBS only     0 arm    0.6912   6054.91 4585.77
   m6gd.4xlarge       64.0     16    1 x 950 NVMe   950 arm    0.7232   6335.23 4650.68
    z1d.2xlarge       64.0      8    1 x 300 NVMe   300 int    0.7440   6517.44 4739.16
     g3s.xlarge       30.5      4        EBS only     0 int    0.7500   6570.00 5553.84
   g4dn.2xlarge       32.0      8     225 GB NVMe     0 int    0.7520   6587.52 4774.20
     m5.4xlarge       64.0     16        EBS only     0 int    0.7680   6727.68 4940.64
    c5d.4xlarge       32.0     16    1 x 400 NVMe   400 int    0.7680   6727.68 4879.32
     c4.4xlarge       30.0     16        EBS only     0 int    0.7960   6972.96 5080.80
     m4.4xlarge       64.0     16        EBS only     0 int    0.8000   7008.00 4992.32
    r6g.4xlarge      128.0     16        EBS only     0 arm    0.8064   7064.06 5121.97
   m5ad.4xlarge       64.0     16    2 x 300 NVMe   600 amd    0.8240   7218.24 5308.56
     x1e.xlarge      122.0      4         1 x 120   120 int    0.8340   7305.84 5177.16
     c3.4xlarge       30.0     16         2 x 160   320 int    0.8400   7358.40 5352.36
      i2.xlarge       30.5      4         1 x 800   800 int    0.8530   7472.28 3836.88
    c5n.4xlarge       42.0     16        EBS only     0 int    0.8640   7568.64 5483.76
   g4ad.4xlarge       64.0     16     600 GB NVMe     0 amd    0.8670   7594.92 5740.08
      p2.xlarge       61.0      4        EBS only     0 int    0.9000   7884.00 6184.56
    r5a.4xlarge      128.0     16        EBS only     0 amd    0.9040   7919.04 5755.32
   i3en.2xlarge       64.0      8   2 x 2500 NVMe  5000 int    0.9040   7919.04 6245.88
    m5d.4xlarge       64.0     16    2 x 300 NVMe   600 int    0.9040   7919.04 5799.12
   r6gd.4xlarge      128.0     16    1 x 950 NVMe   950 arm    0.9216   8073.22 5886.72
     h1.4xlarge       64.0     16    2 x 2000 HDD  4000 int    0.9360   8199.36 6412.32
    m5n.4xlarge       64.0     16        EBS only     0 int    0.9520   8339.52 6088.20
     m2.4xlarge       68.4      8         2 x 840  1680 int    0.9800   8584.80 0.00
   m5zn.3xlarge       48.0     12        EBS only     0 int    0.9910   8681.16 6562.99
     d3.2xlarge       64.0      8    6 x 2000 HDD 12000 int    0.9990   8751.24 6613.80
     r5.4xlarge      128.0     16        EBS only     0 int    1.0080   8830.08 6394.80
   r5ad.4xlarge      128.0     16    2 x 300 NVMe   600 amd    1.0480   9180.48 6683.88
   d3en.2xlarge       32.0      8   4 x 14000 HDD 56000 int    1.0510   9206.76 6964.20
     r4.4xlarge      122.0     16        EBS only     0 int    1.0640   9320.64 6769.73
   m5dn.4xlarge       64.0     16    2 x 300 NVMe   600 int    1.0880   9530.88 6964.20
    c6g.8xlarge       64.0     32        EBS only     0 arm    1.0880   9530.88 6945.80
    z1d.3xlarge       96.0     12    1 x 450 NVMe   450 int    1.1160   9776.16 7104.36
     g3.4xlarge      122.0     16        EBS only     0 int    1.1400   9986.40 7838.45
    r5d.4xlarge      128.0     16    2 x 300 NVMe   600 int    1.1520  10091.52 7323.36
    r5n.4xlarge      128.0     16        EBS only     0 int    1.1920  10441.92 7568.64
    r5b.4xlarge      128.0     16        EBS only     0 int    1.1920  10441.92 7565.17
   g4dn.4xlarge       64.0     16     225 GB NVMe     0 int    1.2040  10547.04 7638.72
   c6gd.8xlarge       64.0     32   1 x 1900 NVMe  1900 arm    1.2288  10764.29 7888.38
    m6g.8xlarge      128.0     32        EBS only     0 arm    1.2320  10792.32 7946.20
    c5a.8xlarge       64.0     32        EBS only     0 amd    1.2320  10792.32 7840.20
     i3.4xlarge      122.0     16   2 x 1900 NVMe  3800 int    1.2480  10932.48 8628.60
     r3.4xlarge      122.0     16         1 x 320   320 int    1.3300  11650.80 7594.92
   r5dn.4xlarge      128.0     16    2 x 300 NVMe   600 int    1.3360  11703.36 8479.68
   i3en.3xlarge       96.0     12   1 x 7500 NVMe  7500 int    1.3560  11878.56 9364.44
   c5ad.8xlarge       64.0     32    2 x 600 NVMe  1200 amd    1.3760  12053.76 8716.20
    m5a.8xlarge      128.0     32        EBS only     0 amd    1.3760  12053.76 8821.32
     d2.2xlarge       61.0      8    6 x 2000 HDD 12000 int    1.3800  12088.80 7297.08
   c6gn.8xlarge       64.0     32        EBS only     0 arm    1.3824  12109.82 9171.54
   m6gd.8xlarge      128.0     32   1 x 1900 NVMe  1900 arm    1.4464  12670.46 9300.49
     c5.9xlarge       72.0     36        EBS only     0 int    1.5300  13402.80 9706.08
     m5.8xlarge      128.0     32        EBS only     0 int    1.5360  13455.36 9881.28
     c4.8xlarge       60.0     36        EBS only     0 int    1.5910  13937.16 10152.84
    r6g.8xlarge      256.0     32        EBS only     0 arm    1.6128  14128.13 10243.07
   c6g.12xlarge       96.0     48        EBS only     0 arm    1.6320  14296.32 10419.14
   m5ad.8xlarge      128.0     32    2 x 600 NVMe  1200 amd    1.6480  14436.48 10608.36
     f1.2xlarge      122.0      8    1 x 470 NVMe   470 int    1.6500  14454.00 11239.08
    x1e.2xlarge      244.0      8         1 x 240   240 int    1.6680  14611.68 10354.32
     c3.8xlarge       60.0     32         2 x 320   640 int    1.6800  14716.80 10695.96
     i2.2xlarge       61.0      8         2 x 800  1600 int    1.7050  14935.80 7673.76
    c5d.9xlarge       72.0     36    1 x 900 NVMe   900 int    1.7280  15137.28 10976.28
   g4ad.8xlarge      128.0     32    1200 GB NVMe     0 amd    1.7340  15189.84 11480.16
    r5a.8xlarge      256.0     32        EBS only     0 amd    1.8080  15838.08 11510.64
    m5d.8xlarge      128.0     32    2 x 600 NVMe  1200 int    1.8080  15838.08 11598.24
   r6gd.8xlarge      256.0     32   1 x 1900 NVMe  1900 arm    1.8432  16146.43 11772.56
  c6gd.12xlarge       96.0     48   2 x 1425 NVMe  2850 arm    1.8432  16146.43 11832.13
   c5a.12xlarge       96.0     48        EBS only     0 amd    1.8480  16188.48 11764.68
   m6g.12xlarge      192.0     48        EBS only     0 arm    1.8480  16188.48 11918.86
     h1.8xlarge      128.0     32    4 x 2000 HDD  8000 int    1.8720  16398.72 12815.88
    m5n.8xlarge      128.0     32        EBS only     0 int    1.9040  16679.04 12167.64
   inf1.6xlarge       48.0     24        EBS only     0 int    1.9040  16679.04 12080.04
    c5n.9xlarge       96.0     36        EBS only     0 int    1.9440  17029.44 12342.84
   m5zn.6xlarge       96.0     24        EBS only     0 int    1.9820  17362.32 13125.98
     d3.4xlarge      128.0     16   12 x 2000 HDD 24000 int    1.9980  17502.48 13227.60
    m4.10xlarge      160.0     40        EBS only     0 int    2.0000  17520.00 12482.12
    cc2.8xlarge       60.5     32         4 x 840  3360 int    2.0000  17520.00 0.00
     r5.8xlarge      256.0     32        EBS only     0 int    2.0160  17660.16 12798.36
    c5.12xlarge       96.0     48        EBS only     0 int    2.0400  17870.40 12947.28
   m5a.12xlarge      192.0     48        EBS only     0 amd    2.0640  18080.64 13236.36
  c5ad.12xlarge       96.0     48    2 x 900 NVMe  1800 amd    2.0640  18080.64 13078.68
  c6gn.12xlarge       96.0     48        EBS only     0 arm    2.0736  18164.74 13757.32
   r5ad.8xlarge      256.0     32    2 x 600 NVMe  1200 amd    2.0960  18360.96 13359.00
   d3en.4xlarge       64.0     16   8 x 14000 HDD 112000 int    2.1030  18422.28 13928.40
     r4.8xlarge      244.0     32        EBS only     0 int    2.1280  18641.28 13539.46
  m6gd.12xlarge      192.0     48   2 x 1425 NVMe  2850 arm    2.1696  19005.70 13951.18
   g4dn.8xlarge      128.0     32     900 GB NVMe     0 int    2.1760  19061.76 13814.52
      c6g.metal      128.0     64        EBS only     0 arm    2.1760  19061.76 13892.48
   c6g.16xlarge      128.0     64        EBS only     0 arm    2.1760  19061.76 13892.48
   m5dn.8xlarge      128.0     32    2 x 600 NVMe  1200 int    2.1760  19061.76 13928.40
    z1d.6xlarge      192.0     24    1 x 900 NVMe   900 int    2.2320  19552.32 14208.72
     g3.8xlarge      244.0     32        EBS only     0 int    2.2800  19972.80 15676.02
    r5d.8xlarge      256.0     32    2 x 600 NVMe  1200 int    2.3040  20183.04 14646.72
    m5.12xlarge      192.0     48        EBS only     0 int    2.3040  20183.04 14830.68
   c5d.12xlarge       96.0     48    2 x 900 NVMe  1800 int    2.3040  20183.04 14637.96
    r5n.8xlarge      256.0     32        EBS only     0 int    2.3840  20883.84 15128.52
    r5b.8xlarge      256.0     32        EBS only     0 int    2.3840  20883.84 15130.34
   r6g.12xlarge      384.0     48        EBS only     0 arm    2.4192  21192.19 15365.04
     c6gd.metal      128.0     64   2 x 1900 NVMe  3800 arm    2.4576  21528.58 15775.88
  c6gd.16xlarge      128.0     64   2 x 1900 NVMe  3800 arm    2.4576  21528.58 15775.88
   c5a.16xlarge      128.0     64        EBS only     0 amd    2.4640  21584.64 15680.40
   m6g.16xlarge      256.0     64        EBS only     0 arm    2.4640  21584.64 15891.52
      m6g.metal      256.0     64        EBS only     0 arm    2.4640  21584.64 15891.52
  m5ad.12xlarge      192.0     48    2 x 900 NVMe  1800 amd    2.4720  21654.72 15916.92
     i3.8xlarge      244.0     32   4 x 1900 NVMe  7600 int    2.4960  21864.96 17265.96
     g2.8xlarge       60.0     32         2 x 120   240 int    2.6000  22776.00 17143.32
     r3.8xlarge      244.0     32         2 x 320   640 int    2.6600  23301.60 15181.08
   r5dn.8xlarge      256.0     32    2 x 600 NVMe  1200 int    2.6720  23406.72 16959.36
   r5a.12xlarge      384.0     48        EBS only     0 amd    2.7120  23757.12 17265.96
   m5d.12xlarge      192.0     48    2 x 900 NVMe  1800 int    2.7120  23757.12 17397.36
   i3en.6xlarge      192.0     24   2 x 7500 NVMe 15000 int    2.7120  23757.12 18737.64
   m5a.16xlarge      256.0     64        EBS only     0 amd    2.7520  24107.52 17642.64
  c5ad.16xlarge      128.0     64   2 x 1200 NVMe  2400 amd    2.7520  24107.52 17441.16
     d2.4xlarge      122.0     16   12 x 2000 HDD 24000 int    2.7600  24177.60 14594.16
  r6gd.12xlarge      384.0     48   2 x 1425 NVMe  2850 arm    2.7648  24219.65 17659.28
  c6gn.16xlarge      128.0     64        EBS only     0 arm    2.7648  24219.65 18343.09
   m5n.12xlarge      192.0     48        EBS only     0 int    2.8560  25018.56 18255.84
  m6gd.16xlarge      256.0     64   2 x 1900 NVMe  3800 arm    2.8928  25340.93 18601.86
     m6gd.metal      256.0     64   2 x 1900 NVMe  3800 arm    2.8928  25340.93 18601.86
    r5.12xlarge      384.0     48        EBS only     0 int    3.0240  26490.24 19193.16
     p3.2xlarge       61.0      8        EBS only     0 int    3.0600  26805.60 21050.28
    c5.18xlarge      144.0     72        EBS only     0 int    3.0600  26805.60 19420.92
    m5.16xlarge      256.0     64        EBS only     0 int    3.0720  26910.72 19771.32
  r5ad.12xlarge      384.0     48    2 x 900 NVMe  1800 amd    3.1440  27541.44 20042.88
   d3en.6xlarge       96.0     24  12 x 14000 HDD 168000 int    3.1540  27629.04 20892.60
    m4.16xlarge      256.0     64        EBS only     0 int    3.2000  28032.00 19971.05
      r6g.metal      512.0     64        EBS only     0 arm    3.2256  28256.26 20486.14
   r6g.16xlarge      512.0     64        EBS only     0 arm    3.2256  28256.26 20486.14
  m5dn.12xlarge      192.0     48 2 x 900 GB NVMe  1800 int    3.2640  28592.64 20892.60
  m5ad.16xlarge      256.0     64    4 x 600 NVMe  2400 amd    3.2960  28872.96 21225.48
     f1.4xlarge      244.0     16    1 x 940 NVMe   940 int    3.3000  28908.00 22478.16
    x1e.4xlarge      488.0     16         1 x 480   480 int    3.3360  29223.36 20708.64
     i2.4xlarge      122.0     16         4 x 800  3200 int    3.4100  29871.60 15347.52
   c5d.18xlarge      144.0     72    2 x 900 NVMe  1800 int    3.4560  30274.56 21961.32
   r5d.12xlarge      384.0     48    2 x 900 NVMe  1800 int    3.4560  30274.56 21970.08
  g4ad.16xlarge      256.0     64    2400 GB NVMe     0 amd    3.4680  30379.68 22960.31
    cr1.8xlarge      244.0     32         2 x 120   240 int    3.5000  30660.00 0.00
   r5b.12xlarge      384.0     48        EBS only     0 int    3.5760  31325.76 22695.51
   r5n.12xlarge      384.0     48        EBS only     0 int    3.5760  31325.76 22697.16
   r5a.16xlarge      512.0     64        EBS only     0 amd    3.6160  31676.16 23012.52
   m5d.16xlarge      256.0     64    4 x 600 NVMe  2400 int    3.6160  31676.16 23196.48
  r6gd.16xlarge      512.0     64   2 x 1900 NVMe  3800 arm    3.6864  32292.86 23545.13
     r6gd.metal      512.0     64   2 x 1900 NVMe  3800 arm    3.6864  32292.86 23545.13
   c5a.24xlarge      192.0     96        EBS only     0 amd    3.6960  32376.96 23529.36
    h1.16xlarge      256.0     64    8 x 2000 HDD 16000 int    3.7440  32797.44 25631.76
   m5n.16xlarge      256.0     64        EBS only     0 int    3.8080  33358.08 24344.04
   c5n.18xlarge      192.0     72        EBS only     0 int    3.8880  34058.88 24676.92
      c5n.metal      192.0     72        EBS only     0 int    3.8880  34058.88 24676.92
  g4dn.12xlarge      192.0     48     900 GB NVMe     0 int    3.9120  34269.12 24825.84
     m5zn.metal      192.0     48        EBS only     0 int    3.9641  34725.52 26251.97
  m5zn.12xlarge      192.0     48        EBS only     0 int    3.9641  34725.52 26251.97
     d3.8xlarge      256.0     32   24 x 2000 HDD 48000 int    3.9955  35000.76 26460.54
  r5dn.12xlarge      384.0     48 2 x 900 GB NVMe  1800 int    4.0080  35110.08 25439.04
    r5.16xlarge      512.0     64        EBS only     0 int    4.0320  35320.32 25587.96
       c5.metal      192.0     96        EBS only     0 int    4.0800  35740.80 25894.56
    c5.24xlarge      192.0     96        EBS only     0 int    4.0800  35740.80 25894.56
   m5a.24xlarge      384.0     96        EBS only     0 amd    4.1280  36161.28 26463.96
  c5ad.24xlarge      192.0     96   2 x 1900 NVMe  3800 amd    4.1280  36161.28 26157.36
  r5ad.16xlarge      512.0     64    4 x 600 NVMe  2400 amd    4.1920  36721.92 26718.00
   d3en.8xlarge      128.0     32  16 x 14000 HDD 224000 int    4.2058  36842.46 27852.86
    r4.16xlarge      488.0     64        EBS only     0 int    4.2560  37282.56 27078.91
  m5dn.16xlarge      256.0     64    4 x 600 NVMe  2400 int    4.3520  38123.52 27856.80
  g4dn.16xlarge      256.0     64     900 GB NVMe     0 int    4.3520  38123.52 27620.28
      z1d.metal      384.0     48    2 x 900 NVMe  1800 int    4.4640  39104.64 28417.44
   z1d.12xlarge      384.0     48    2 x 900 NVMe  1800 int    4.4640  39104.64 28417.44
    g3.16xlarge      488.0     64        EBS only     0 int    4.5600  39945.60 31352.04
    hs1.8xlarge      117.0     16   24 x 2000 HDD 48000 int    4.6000  40296.00 0.00
       m5.metal      384.0     96        EBS only     0 int    4.6080  40366.08 29652.60
   r5d.16xlarge      512.0     64    4 x 600 NVMe  2400 int    4.6080  40366.08 29293.44
    m5.24xlarge      384.0     96        EBS only     0 int    4.6080  40366.08 29652.60
      c5d.metal      192.0     96    4 x 900 NVMe  3600 int    4.6080  40366.08 29275.92
   c5d.24xlarge      192.0     96    4 x 900 NVMe  3600 int    4.6080  40366.08 29275.92
   r5b.16xlarge      512.0     64        EBS only     0 int    4.7680  41767.68 30260.68
   r5n.16xlarge      512.0     64        EBS only     0 int    4.7680  41767.68 30257.04
  m5ad.24xlarge      384.0     96    4 x 900 NVMe  3600 amd    4.9440  43309.44 31833.84
    i3.16xlarge      488.0     64   8 x 1900 NVMe 15200 int    4.9920  43729.92 34523.16
       i3.metal      512.0     64   8 x 1900 NVMe 15200 int    4.9920  43729.92 34523.16
  r5dn.16xlarge      512.0     64    4 x 600 NVMe  2400 int    5.3440  46813.44 33918.72
   r5a.24xlarge      768.0     96        EBS only     0 amd    5.4240  47514.24 34523.16
      m5d.metal      384.0     96    4 x 900 NVMe  3600 int    5.4240  47514.24 34794.72
   m5d.24xlarge      384.0     96    4 x 900 NVMe  3600 int    5.4240  47514.24 34794.72
  i3en.12xlarge      384.0     48   4 x 7500 NVMe 30000 int    5.4240  47514.24 37475.28
     d2.8xlarge      244.0     36   24 x 2000 HDD 48000 int    5.5200  48355.20 29197.08
   m5n.24xlarge      384.0     96        EBS only     0 int    5.7120  50037.12 36511.68
    r5.24xlarge      768.0     96        EBS only     0 int    6.0480  52980.48 38386.32
       r5.metal      768.0     96        EBS only     0 int    6.0480  52980.48 38386.32
  r5ad.24xlarge      768.0     96    4 x 900 NVMe  3600 amd    6.2880  55082.88 40077.00
  d3en.12xlarge      192.0     48  24 x 14000 HDD 336000 int    6.3086  55263.69 41779.33
  m5dn.24xlarge      384.0     96    4 x 900 NVMe  3600 int    6.5280  57185.28 41793.96
    x1.16xlarge      976.0     64        1 x 1920  1920 int    6.6690  58420.44 41399.76
    x1e.8xlarge      976.0     32         1 x 960   960 int    6.6720  58446.72 41417.28
     i2.8xlarge      244.0     32         8 x 800  6400 int    6.8200  59743.20 30703.80
      r5d.metal      768.0     96    4 x 900 NVMe  3600 int    6.9120  60549.12 43948.92
   r5d.24xlarge      768.0     96    4 x 900 NVMe  3600 int    6.9120  60549.12 43948.92
      r5b.metal      768.0     96        EBS only     0 int    7.1520  62651.52 45391.03
   r5b.24xlarge      768.0     96        EBS only     0 int    7.1520  62651.52 45391.03
   r5n.24xlarge      768.0     96        EBS only     0 int    7.1520  62651.52 45394.32
     p2.8xlarge      488.0     32        EBS only     0 int    7.2000  63072.00 49502.76
  inf1.24xlarge      192.0     96        EBS only     0 int    7.6150  66707.40 48328.92
     g4dn.metal      384.0     96 2 x 900 GB NVMe  1800 int    7.8240  68538.24 49651.68
  r5dn.24xlarge      768.0     96    4 x 900 NVMe  3600 int    8.0160  70220.16 50878.08
     i3en.metal      768.0     96   8 x 7500 NVMe 60000 int   10.8480  95028.48 74950.56
  i3en.24xlarge      768.0     96   8 x 7500 NVMe 60000 int   10.8480  95028.48 74950.56
     p3.8xlarge      244.0     32        EBS only     0 int   12.2400 107222.40 84209.88
    f1.16xlarge      976.0     64    4 x 940 NVMe  3760 int   13.2000 115632.00 89909.14
    x1.32xlarge     1952.0    128        2 x 1920  3840 int   13.3380 116840.88 82799.52
   x1e.16xlarge     1952.0     64        1 x 1920  1920 int   13.3440 116893.44 82843.32
    p2.16xlarge      732.0     64        EBS only     0 int   14.4000 126144.00 99005.52
    p3.16xlarge      488.0     64        EBS only     0 int   24.4800 214444.80 168428.52
   x1e.32xlarge     3904.0    128        2 x 1920  3840 int   26.6880 233786.88 165677.88
  p3dn.24xlarge      768.0     96    2 x 900 NVMe  1800 int   31.2120 273417.12 193552.20
   p4d.24xlarge     1152.0     96        8 x 1000  8000 int   32.7726 287087.98 212084.84

# list family types and their general usecase
$ ./ec2price -family
   m1 2006       main
   c1 2008        cpu
   m2 2009       main
  cc1 2010 cluster-co
   t1 2010      burst
  cg1 2010        gpu
  cc2 2011 cluster-co
  hi1 2012        ssd
   m3 2012       main
  hs1 2012  dense-hdd
  cr1 2013 cluster-co
   c3 2013        cpu
   g2 2013        gpu
   i2 2013        ssd
   r3 2014   more-mem
   t2 2014      burst
   c4 2015        cpu
   d2 2015  dense-hdd
   m4 2015       main
   x1 2016 mem-xtreme
   p2 2016        gpu
   f1 2016       fpga
   r4 2016   more-mem
   i3 2016        ssd
   c5 2016        cpu
   g3 2017        gpu
  x1e 2017 mem-xtreme
   p3 2017        gpu
   m5 2017       main
   h1 2017  dense-hdd
  c5d 2018        cpu nvme
  m5d 2018       main nvme
  z1d 2018  high-freq
   r5 2018   more-mem
   t3 2018      burst
  g3s 2018        gpu
  m5a 2018       main amd
  r5a 2018   more-mem amd
  c5n 2018        cpu net
   a1 2018        arm
 p3dn 2018        gpu nvme,net
   g4 2019        gpu
 m5ad 2019       main amd,nvme
  r5d 2019   more-mem nvme
 r5ad 2019   more-mem amd,nvme
 i3en 2019        ssd net
 g4dn 2019        gpu gpu-amd
 r5dn 2019   more-mem nvme,net
  r5n 2019   more-mem net
 m5dn 2019       main nvme,net
  m5n 2019       main net
 inf1 2019  inference
  t3a 2019      burst amd
  c5a 2020        cpu amd
 c5ad 2020        cpu amd,nvme
  c6g 2020        cpu graviton
 c6gn 2020        cpu graviton,net
 c6gd 2020        cpu graviton,nvme
   d3 2020  dense-hdd
 d3en 2020  dense-hdd net
 g4ad 2020        gpu amd,nvme
 m5zn 2020       main net,high-freq
  m6g 2020       main graviton
 m6gn 2020       main graviton,net
  p4d 2020        gpu nvme
  r5b 2020   more-mem ebs-optimized
  r6g 2020   more-mem graviton
 r6gd 2020   more-mem graviton,nvme
  t4g 2020      burst graviton
 x2gd 2021 mem-xtreme graviton,nvme


# flags
$ ./ec2price -help
Usage of ./ec2price:
  -check_family
        Check family types against instance types (for missing types)
  -csv
        output as csv
  -family
        Print family type information
  -fetch_offers
        Fetch offers and price file to disk
  -region string
        AWS Region (default "us-east-1")

```

## License

MIT
