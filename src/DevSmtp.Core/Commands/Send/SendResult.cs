namespace DevSmtp.Core.Commands
{
    public sealed class SendResult : CommandResult
    {
        public SendResult()
        {
        }

        public SendResult(Exception error)
            : base(error)
        {
        }
    }
}
