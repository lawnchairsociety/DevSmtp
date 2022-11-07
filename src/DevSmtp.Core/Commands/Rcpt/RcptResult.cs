namespace DevSmtp.Core.Commands
{
    public sealed class RcptResult : CommandResult
    {
        public RcptResult()
        {
        }

        public RcptResult(Exception error)
            : base(error)
        {
        }
    }
}
