Feature: renew action with `-csr service` option

  As a user
  I want to renew certificates that were enrolled by the app
  Using `-csr service` option meaning that new private key and CSR are generated on service side

  It requires key password typed interactively or -key-password option to be used to download key from TPP

  - for TPP:
    - certificate is requested to be renewed on service side.
      if it's "User Provided CSR",
        then "Waiting for new CSR" error returned,

      if it's "Service Generated CSR", then
        if policy allows key reuse, then old CSR is signed, or new key & CSR generated otherwise

      if it's "Service Generated CSR" and PKCS#12 format is specified, then
        if policy allows key reuse, then old CSR is signed, or new key & CSR generated otherwise and it should return PKCS#12 file

  - for Condor:
    - certificate is requested to be renewed on service side
      if policy allows key reuse, then old CSR is signed, error returns otherwise

  Background:
    And the default aruba exit timeout is 180 seconds

  Scenario: where it should return an error if renew is used in TPP with -csr=service and empty -key-password
    When I renew the certificate in TPP with flags -id xxx -no-prompt -csr service
    Then it should fail with "-key-password cannot be empty in -csr service mode for TPP unless -no-pickup specified"

  Scenario: renew user-provided-CSR certificate in TPP with `-csr service` option
    Given I enroll random certificate using TPP with -no-prompt -key-file k.pem -cert-file c.pem
      And it should write private key to the file "k.pem"
      And it should write certificate to the file "c.pem"
      And it should output Pickup ID
    When I renew the certificate in TPP using the same Pickup ID with flags -no-prompt -cert-file c1.pem -key-file k1.pem -csr service -key-password Passcode123!
    Then it should fail with "Status: 400"

  Scenario: renew service-generated-CSR certificate in TPP with `-csr service` option
    Given I enroll random certificate using TPP with -csr service -key-file k.pem -cert-file c.pem -key-password Passcode123!
      And it should write private key to the file "k.pem"
      And it should write certificate to the file "c.pem"
      And it should output Pickup ID
    When I renew the certificate in TPP using the same Pickup ID with flags -cert-file c1.pem -key-file k1.pem -csr service -key-password Passcode123!
      Then it should retrieve certificate
      And it should write private key to the file "k1.pem"
      And it should write certificate to the file "c1.pem"
      And private key in "k1.pem" with passphrase and certificate in "c1.pem" should have the same modulus

  @TODO # This scenario is not currently working due to reuse csr not being implemented yet in OutagePredict
  Scenario: renew certificate in VaaS with -csr=service which is working only if Zone's policy allows key reuse
    Given I enroll random certificate using VaaS with -csr service -no-prompt -key-file k.pem -cert-file c.pem -timeout 180
      And it should write private key to the file "k.pem"
      And it should write certificate to the file "c.pem"
      And it should output Pickup ID
    Then I renew the certificate in VaaS using the same Pickup ID with flags -csr service -no-prompt -cert-file c1.pem -key-file k1.pem
      And it should retrieve certificate
      But it should not output private key
      And it should write certificate to the file "c1.pem"
      And certificate in "k.pem" and certificate in "k1.pem" should have the same modulus
      And certificate in "c.pem" and certificate in "c1.pem" should not have the same serial

  Scenario: renew service-generated-CSR certificate in TPP with `-csr service` option with PKCS12 flag
    Given I enroll random certificate using TPP with -csr service -key-password Passcode123! -key-file k.pem -cert-file c.pem -key-password Passcode123!
      And it should write private key to the file "k.pem"
      And it should write certificate to the file "c.pem"
      And it should output Pickup ID
    When I renew the certificate in TPP using the same Pickup ID with flags -csr service -key-password Passcode123! -file all.p12 -format pkcs12
      Then it should retrieve certificate
      And "all.p12" should be PKCS#12 archive with password "Passcode123!"
