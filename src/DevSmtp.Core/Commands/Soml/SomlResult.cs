namespace DevSmtp.Core.Commands
{
    public sealed class SomlResult : CommandResult
    {
        public SomlResult()
        {
        }

        public SomlResult(Exception error)
            : base(error)
        {
        }
    }
}
