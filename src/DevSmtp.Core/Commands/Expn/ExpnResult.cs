namespace DevSmtp.Core.Commands
{
    public sealed class ExpnResult : CommandResult
    {
        public ExpnResult()
        {
        }

        public ExpnResult(Exception error)
            : base(error)
        {
        }
    }
}
