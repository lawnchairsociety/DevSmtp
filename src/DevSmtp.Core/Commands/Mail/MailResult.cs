namespace DevSmtp.Core.Commands
{
    public sealed class MailResult : CommandResult
    {
        public MailResult()
        {
        }

        public MailResult(Exception error)
            : base(error)
        {
        }
    }
}
