namespace DevSmtp.Core.Commands
{
    public sealed class RsetResult : CommandResult
    {
        public RsetResult()
        {
        }

        public RsetResult(Exception error)
            : base(error)
        {
        }
    }
}
