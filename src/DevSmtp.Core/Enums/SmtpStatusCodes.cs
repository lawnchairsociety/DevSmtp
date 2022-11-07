namespace DevSmtp.Core.Enums
{
    public enum SmtpStatusCodes
    {
        GeneralFailure = -1,                       // The transaction could not occur. You receive this error when the specified SMTP host cannot be found.
        SystemStatus = 211,                        // A system status or system Help reply.
        HelpMessage = 214,                         // A Help message was returned by the service.
        ServiceReady = 220,                        // The SMTP service is ready.
        ServiceClosingTransmissionChannel = 221,   // The SMTP service is closing the transmission channel.
        Ok = 250,                                  // The email was successfully sent to the SMTP service.
        UserNotLocalWillForward = 251,             // The user mailbox is not located on the receiving server; the server forwards the email.
        CannotVerifyUserWillAttemptDelivery = 252, // The specified user is not local, but the receiving SMTP service accepted the message and attempted to deliver it.

        StartMailInput = 354,                      // The SMTP service is ready to receive the email content.

        ServiceNotAvailable = 421,                 // The SMTP service is not available; the server is closing the transmission channel.
        MailboxBusy = 450,                         // The destination mailbox is in use.
        LocalErrorInProcessing = 451,              // The SMTP service cannot complete the request.
        InsufficientStorage = 452,                 // The SMTP service does not have sufficient storage to complete the request.
        ClientNotPermitted = 454,                  // The client was not authenticated or is not allowed to send mail using the specified SMTP host.

        CommandUnrecognized = 500,                 // The SMTP service does not recognize the specified command.
        SyntaxError = 501,                         // The syntax used to specify a command or parameter is incorrect.
        CommandNotImplemented = 502,               // The SMTP service does not implement the specified command.
        BadCommandSequence = 503,                  // The commands were sent in the incorrect sequence.
        CommandParameterNotImplemented = 504,      // The SMTP service does not implement the specified command parameter.
        MustIssueStartTlsFirst = 530,              // The SMTP server is configured to accept only TLS connections.
        MailboxUnavailable = 550,                  // The destination mailbox was not found or could not be accessed.
        UserNotLocalTryAlternatePath = 551,        // The user mailbox is not located on the receiving server; resend using the supplied address information.
        ExceededStorageAllocation = 552,           // The message is too large to be stored in the destination mailbox.
        MailboxNameNotAllowed = 553,               // The syntax used to specify the destination mailbox is incorrect.
        TransactionFailed = 554,                   // The transaction failed.
    }
}
