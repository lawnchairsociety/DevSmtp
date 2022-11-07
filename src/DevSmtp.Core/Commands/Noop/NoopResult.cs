namespace DevSmtp.Core.Commands
{
    public sealed class NoopResult : CommandResult
    {
        public NoopResult()
        {
        }

        public NoopResult(Exception error)
            : base(error)
        {
        }
    }
}
